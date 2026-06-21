import { env as priv } from '$env/dynamic/private';
import { env as pub } from '$env/dynamic/public';

/** @type {import('./$types').LayoutServerLoad} */
export async function load({ fetch }) {
	const endpoint = `${pub.PUBLIC_INTERNAL_API_HOST}/api/transcriptions`;
	const headers = {};
	if (priv.WHISHPER_API_KEY) headers['X-API-Key'] = priv.WHISHPER_API_KEY;

	let ts = [];
	try {
		const response = await fetch(endpoint, { headers });
		if (response.ok) {
			ts = (await response.json()) ?? [];
		}
	} catch (e) {
		ts = [];
	}

	return {
		transcriptions: ts,
		authEnabled: !!priv.WHISHPER_API_KEY
	};
}
