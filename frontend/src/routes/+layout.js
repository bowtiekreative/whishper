import { transcriptions } from '$lib/stores';

/** @type {import('./$types').LayoutLoad} */
export function load({ data }) {
	const ts = data?.transcriptions;
	transcriptions.set(Array.isArray(ts) && ts.length > 0 ? ts : []);
	return data;
}
