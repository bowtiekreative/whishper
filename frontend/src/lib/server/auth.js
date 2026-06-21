import crypto from 'node:crypto';
import { env } from '$env/dynamic/private';

export const SESSION_COOKIE = 'whishper_session';

// Auth is enabled only when an API key is configured (backwards compatible).
export function authEnabled() {
	return !!env.WHISHPER_API_KEY;
}

// The key the frontend server uses to authenticate its own internal API calls.
export function apiKey() {
	return env.WHISHPER_API_KEY || '';
}

// validateToken mirrors the backend's HMAC session token verification so the
// SvelteKit server can gate pages without an extra round-trip to the backend.
export function validateToken(token) {
	const secret = env.WHISHPER_API_KEY;
	if (!secret || !token) return false;

	const parts = token.split('.');
	if (parts.length !== 2) return false;

	const expected = crypto.createHmac('sha256', secret).update(parts[0]).digest('base64url');
	const a = Buffer.from(expected);
	const b = Buffer.from(parts[1]);
	if (a.length !== b.length || !crypto.timingSafeEqual(a, b)) return false;

	try {
		const payload = JSON.parse(Buffer.from(parts[0], 'base64url').toString());
		if (Math.floor(Date.now() / 1000) > payload.exp) return false;
	} catch {
		return false;
	}
	return true;
}
