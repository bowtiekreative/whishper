import { env as priv } from '$env/dynamic/private';
import { env as pub } from '$env/dynamic/public';

/** @type {import('./$types').PageServerLoad} */
export async function load({ params, fetch }) {
	const endpoint = `${pub.PUBLIC_INTERNAL_API_HOST}/api/transcriptions/${params.id}`;
	const headers = {};
	if (priv.WHISHPER_API_KEY) headers['X-API-Key'] = priv.WHISHPER_API_KEY;

	let transcription = null;
	try {
		const response = await fetch(endpoint, { headers });
		if (response.ok) {
			transcription = await response.json();
		}
	} catch (e) {
		transcription = null;
	}

	return { transcription };
}
