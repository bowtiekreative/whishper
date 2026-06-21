import { redirect } from '@sveltejs/kit';
import { authEnabled, validateToken, SESSION_COOKIE } from '$lib/server/auth';

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
	const { pathname } = event.url;

	// Let static assets (anything with a file extension) and the login page
	// through, and skip the gate entirely when auth is disabled.
	const isAsset = pathname.includes('.');
	const isLogin = pathname === '/login' || pathname.startsWith('/login');

	if (authEnabled() && !isAsset && !isLogin) {
		const token = event.cookies.get(SESSION_COOKIE);
		if (!validateToken(token)) {
			throw redirect(302, '/login');
		}
	}

	return resolve(event);
}
