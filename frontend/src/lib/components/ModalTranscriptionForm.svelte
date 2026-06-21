<script>
	import { validateURL, CLIENT_API_HOST } from '$lib/utils.js';
	import { env } from '$env/dynamic/public';
	import { uploadProgress } from '$lib/stores';

	import toast from 'svelte-french-toast';

	let errorMessage = '';
	let disableSubmit = true;
	let modelSize = 'small';
	let language = 'auto';
	let sourceUrl = '';
	let fileInput;
	let device = env.PUBLIC_WHISHPER_PROFILE == 'gpu' ? 'cuda' : 'cpu';

	let languages = [
		'auto',
		'ar',
		'be',
		'bg',
		'bn',
		'ca',
		'cs',
		'cy',
		'da',
		'de',
		'el',
		'en',
		'es',
		'fr',
		'it',
		'ja',
		'nl',
		'pl',
		'pt',
		'ru',
		'sk',
		'sl',
		'sv',
		'tk',
		'tr',
		'zh'
	];
	let models = [
		'tiny',
		'tiny.en',
		'base',
		'base.en',
		'small',
		'small.en',
		'medium',
		'medium.en',
		'large-v2',
		'large-v3'
	];
	// Sort the languages
	languages.sort((a, b) => {
		if (a == 'auto') return -1;
		if (b == 'auto') return 1;
		return a.localeCompare(b);
	});

	// Uploads a single job (either a file or, when file is null, the sourceUrl).
	// progressCb receives 0-100 for the current item's upload progress.
	function uploadOne(file, progressCb) {
		let formData = new FormData();
		formData.append('language', language);
		formData.append('modelSize', modelSize);
		formData.append('device', device == 'cuda' || device == 'cpu' ? device : 'cpu');

		if (file) {
			formData.append('sourceUrl', '');
			formData.append('file', file);
		} else {
			formData.append('sourceUrl', sourceUrl);
		}

		return new Promise((resolve, reject) => {
			const xhr = new XMLHttpRequest();

			xhr.upload.addEventListener('progress', (event) => {
				if (event.lengthComputable) {
					progressCb(Math.round((event.loaded * 100) / event.total));
				}
			});

			xhr.addEventListener('load', () => {
				if (xhr.status === 200) {
					resolve(xhr.response);
				} else if (xhr.status === 401) {
					reject('unauthorized');
				} else {
					reject(xhr.statusText || 'upload failed');
				}
			});

			xhr.addEventListener('error', () => reject('network error'));

			xhr.open('POST', `${CLIENT_API_HOST}/api/transcriptions`);
			// Send the session cookie so the request is authenticated.
			xhr.withCredentials = true;
			xhr.send(formData);
		});
	}

	// Function that sends one or more jobs to the backend.
	async function sendForm() {
		if (sourceUrl && !validateURL(sourceUrl)) {
			toast.error('You must enter a valid URL.');
			return;
		}

		const files = fileInput && fileInput.files ? Array.from(fileInput.files) : [];

		if (!sourceUrl && files.length === 0) {
			toast.error('No file or URL.');
			return;
		}

		// A source URL takes precedence and is treated as a single job.
		const jobs = sourceUrl ? [null] : files;
		const total = jobs.length;
		let succeeded = 0;

		for (let i = 0; i < total; i++) {
			try {
				await uploadOne(jobs[i], (pct) => {
					// Combine per-file progress into an overall percentage.
					uploadProgress.set(Math.round(((i + pct / 100) / total) * 100));
				});
				succeeded++;
			} catch (err) {
				if (err === 'unauthorized') {
					toast.error('Session expired. Please log in again.');
					window.location.href = '/login';
					return;
				}
				const name = jobs[i] ? jobs[i].name : sourceUrl;
				toast.error(`Failed: ${name}`);
			}
		}

		uploadProgress.set(0);
		if (succeeded > 0) {
			toast.success(
				total > 1 ? `Queued ${succeeded}/${total} files!` : 'Success!'
			);
		}

		// Reset the form inputs.
		sourceUrl = '';
		if (fileInput) fileInput.value = '';
	}

	// Reactive statement
	$: if (sourceUrl && !validateURL(sourceUrl)) {
		errorMessage = 'Enter a valid URL';
		disableSubmit = true;
	} else {
		errorMessage = '';
		disableSubmit = false;
	}
</script>

<dialog id="modalNewTranscription" class="modal">
	<form method="dialog" class="modal-box">
		<button class="absolute btn btn-sm btn-circle btn-ghost right-2 top-2">✕</button>
		{#if errorMessage != ''}
			<div class="alert alert-error">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="w-6 h-6 stroke-current shrink-0"
					fill="none"
					viewBox="0 0 24 24"
					><path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
					/></svg
				>
				<span>{errorMessage}</span>
			</div>
		{/if}
		<div class="mt-0 space-y-2">
			<div class="w-full max-w-xs form-control">
				<label for="file" class="label">
					<span class="label-text">Pick one or more files</span>
				</label>
				<input
					name="file"
					bind:this={fileInput}
					type="file"
					multiple
					class="w-full max-w-xs file-input file-input-sm file-input-bordered file-input-primary"
				/>
			</div>

			<div class="w-full max-w-xs form-control">
				<label for="sourceUrl" class="label">
					<span class="label-text">Or a source URL</span>
				</label>
				<input
					name="sourceUrl"
					bind:value={sourceUrl}
					type="text"
					placeholder="https://youtube.com/watch?v=Hd33fCdW"
					class="w-full max-w-xs input input-sm input-bordered input-primary"
				/>
			</div>
		</div>

		<div class="mb-0 divider" />
		<!-- Whisper Configuration -->
		<div class="flex space-x-4">
			<div class="w-full max-w-xs form-control">
				<label for="modelSize" class="label">
					<span class="label-text">Whisper model</span>
				</label>
				<select name="modelSize" bind:value={modelSize} class="select select-bordered">
					{#each models as m}
						<option value={m}>{m}</option>
					{/each}
				</select>
			</div>

			<div class="w-full max-w-xs form-control">
				<label for="language" class="label">
					<span class="label-text">Language</span>
				</label>
				<select name="language" bind:value={language} class="select select-bordered">
					{#each languages as l}
						<option value={l}>{l}</option>
					{/each}
				</select>
			</div>

			<div class="w-full max-w-xs form-control">
				<label for="language" class="label">
					<span class="label-text">Device</span>
				</label>
				<select name="device" bind:value={device} class="select select-bordered">
					{#if env.PUBLIC_WHISHPER_PROFILE == 'gpu'}
						<option selected value="cuda">GPU</option>
						<option value="cpu">CPU</option>
					{:else}
						<option selected value="cpu">CPU</option>
						<option disabled value="cuda">GPU</option>
					{/if}
				</select>
			</div>
		</div>

		<div class="mb-0 divider" />
		<!--Actions-->
		<button class="btn btn-wide btn-primary" on:click={sendForm} disabled={disableSubmit}
			>Start</button
		>
	</form>
	<form method="dialog" class="modal-backdrop">
		<button>close</button>
	</form>
</dialog>
