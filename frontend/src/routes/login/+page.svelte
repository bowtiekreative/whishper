<script>
	import { goto } from '$app/navigation';
	import toast, { Toaster } from 'svelte-french-toast';

	let username = '';
	let password = '';
	let loading = false;

	async function login() {
		loading = true;
		try {
			const res = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password })
			});

			if (res.ok) {
				toast.success('Welcome back!');
				// Full reload so the new session cookie is picked up by the server hook.
				window.location.href = '/';
			} else if (res.status === 401) {
				toast.error('Invalid credentials');
			} else {
				toast.error('Login failed');
			}
		} catch (e) {
			toast.error('Could not reach the server');
		} finally {
			loading = false;
		}
	}
</script>

<Toaster />

<main class="flex items-center justify-center min-h-screen px-4">
	<form
		class="w-full max-w-sm p-8 space-y-4 card bg-neutral text-neutral-content"
		on:submit|preventDefault={login}
	>
		<h1 class="flex items-center justify-center space-x-3 text-3xl font-bold">
			<img class="w-12 h-12" src="/logo.svg" alt="Whishper logo" />
			<span>Whishper</span>
		</h1>
		<p class="text-center opacity-70">Please sign in to continue</p>

		<div class="w-full form-control">
			<label for="username" class="label"><span class="label-text">Email</span></label>
			<input
				id="username"
				type="text"
				bind:value={username}
				autocomplete="username"
				class="w-full input input-bordered input-primary"
				placeholder="you@example.com"
			/>
		</div>

		<div class="w-full form-control">
			<label for="password" class="label"><span class="label-text">Password</span></label>
			<input
				id="password"
				type="password"
				bind:value={password}
				autocomplete="current-password"
				class="w-full input input-bordered input-primary"
				placeholder="••••••••"
			/>
		</div>

		<button type="submit" class="w-full btn btn-primary" disabled={loading}>
			{loading ? 'Signing in…' : 'Sign in'}
		</button>
	</form>
</main>
