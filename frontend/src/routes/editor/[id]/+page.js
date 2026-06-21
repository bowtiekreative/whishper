import { currentTranscription } from '$lib/stores';

/** @type {import('./$types').PageLoad} */
export function load({ data }) {
	currentTranscription.set(data?.transcription);
	return data;
}
