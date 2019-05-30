class MemoryStorage {
	constructor () {
		this.memory = {};
	}

	getItem (key) {
		if (key in this.memory) {
			return this.memory[key];
		} else {
			return null;
		}
	}

	setItem (key, value) {
		this.memory[key] = value;
	}

	removeItem (key) {
		delete this.memory[key];
	}

	clear () {
		this.memory = {};
	}
}
const memoryStorage = new MemoryStorage();

function testStorage (storage) {
	if (storage) {
		try {
			const key = '__testing-storage__';
			storage.setItem(key, key);
			storage.removeItem(key);
			return true;
		} catch (err) {}
	}
	return false;
}

function getStorage () {
	const storages = [
		window.localStorage,
		window.sessionStorage,
		memoryStorage,
	];
	for (const s of storages) {
		if (testStorage(s)) {
			return s;
		}
	}
	return null;
}
const storage = getStorage();

export function getItem (key) {
	if (!storage) {
		throw new Error('cache unavailable');
	}
	return storage.getItem(key);
}

export function hasItem (key) {
	if (!storage) {
		throw new Error('cache unavailable');
	}
	return storage.getItem(key) !== null;
}

export function setItem (key, value) {
	if (!storage) {
		throw new Error('cache unavailable');
	}
	storage.setItem(key, value);
}

export function removeItem (key) {
	if (!storage) {
		throw new Error('cache unavailable');
	}
	storage.removeItem(key);
}

export default {
	getItem,
	hasItem,
	setItem,
	removeItem,
}
