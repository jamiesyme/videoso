import Vue from 'vue';
import CacheUtils from '@/utils/cache';

class Credentials {
	constructor (email, password) {
		this.email = email;
		this.password = password;
	}

	match (creds) {
		function cleanEmail (email) {
			return email.toLowerCase().trim();
		}
		function cleanPassword (password) {
			return password;
		}
		const e1 = cleanEmail(this.email);
		const e2 = cleanEmail(creds.email);
		const p1 = cleanPassword(this.password);
		const p2 = cleanPassword(creds.password);
		return e1 === e2 && p1 === p2;
	}
}

const state = Vue.observable({
	masterCreds: new Credentials(
		process.env.VUE_APP_LOGIN_EMAIL,
		process.env.VUE_APP_LOGIN_PASSWORD,
	),
	hasInited: false,
	currentCreds: null,
});

function saveCredentials (creds) {
	const str = JSON.stringify(creds);
	CacheUtils.setItem('credentials', str);
}

function loadCredentials () {
	if (CacheUtils.hasItem('credentials')) {
		const str = CacheUtils.getItem('credentials');
		const obj = JSON.parse(str);
		return new Credentials(obj.email, obj.password);
	} else {
		return null;
	}
}

function clearCredentials () {
	if (CacheUtils.hasItem('credentials')) {
		CacheUtils.removeItem('credentials');
	}
}

/**
 * Loads previously cached credentials and compares them against the most recent
 * master credentials, logging out if necessary. Idempotent.
 */
export function init () {
	if (!state.hasInited) {
		state.currentCreds = loadCredentials();
		if (state.currentCreds) {
			const current = state.currentCreds;
			const master = state.masterCreds;
			const stillValid = master.match(current);
			if (!stillValid) {
				logout();
			}
		}
		state.hasInited = true;
	}
}

/**
 * Compares given credentials against the master credentials. If they match, the
 * credentials are cached and `true` is returned. Otherwise, any previously
 * cached credentials are cleared and `false` is returned.
 *
 * @param {string} email
 * @param {string} password
 * @returns {boolean}
 */
export function login (email, password) {
	const creds = new Credentials(email, password);
	if (state.masterCreds.match(creds)) {
		state.currentCreds = creds;
		saveCredentials(creds);
		return true;
	} else {
		state.currentCreds = null;
		clearCredentials();
		return false;
	}
}

/**
 * Clears previously cached credentials.
 */
export function logout () {
	state.currentCreds = null;
	clearCredentials();
}

/**
 * Returns true if currently logged in. Calls `init` automatically.
 *
 * @returns {boolean}
 */
export function loggedIn () {
	init();
	return !!state.currentCreds;
}

export default {
	init,
	login,
	logout,
	loggedIn,
}
