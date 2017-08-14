/**
 * Public functions:
 *   login(emailAddress, password)
 *   logout()
 *   isLoggedIn()
 *   authenticatedAjax(url, data)
 *
 * Private functions:
 *   getRefreshToken(emailAddress, password)
 *   createRefreshToken(emailAddress, password)
 *   destroyRefreshToken(refreshTokenId)
 *   getAccessToken(refreshTokenId)
 *   releaseAccessToken()
 *   isAccessTokenValid(accessToken)
 *   generateAccessToken(refreshTokenId)
 *
 * Storage variables:
 *   localStorage['refreshTokenId']
 *   sessionStorage['accessToken']
 *   sessionStorage['loggedIn']
 */

/**
 * @param {string} name
 * @param {string} [defaultMessage]
 * @return {Error}
 */
function makeCustomError(name, defaultMessage) {
	var e = function(message) {
		this.name = name;
		this.message = message || defaultMessage;
		this.stack = (new Error()).stack;
	};
	e.prototype = Object.create(Error.prototype);
	e.prototype.constructor = e;
	return e;
}

var InvalidCredentialsError = makeCustomError('InvalidCredentialsError');

var LoggedInError = makeCustomError('LoggedInError');

var LoggedOutError = makeCustomError('LoggedOutError');

var NotFoundError = makeCustomError('NotFoundError');

var ServerError = makeCustomError('ServerError');

/**
 * @promise LoginPromise
 * @reject {InvalidCredentialsError} 
 * @reject {LoggedInError} Must log out before logging in again.
 * @reject {ServerError}
 */

/**
 * Creates a refresh token using the provided credentials. The token is stored in
 * long-term storage so that the user is only asked to log in once.
 *
 * @param {string} emailAddress
 * @param {string} password
 * @return {LoginPromise}
 * 
 */
function login(emailAddress, password) {
	if (!emailAddress || !password) {
		throw new TypeError();
	}
	return new Promise(function(resolve, reject) {

		// If we are already logged in, we shouldn't be able to log in again with
		// (potentially) different credentials. Instead, the user should logout
		// first.
		if (isLoggedIn()) {
			reject(new LoggedInError());
			return;
		}

		// This will generate and cache a new refresh token
		getRefreshToken(emailAddress, password)
			.then(function() {
				window.sessionStorage.setItem('loggedIn', 'true');
				resolve();
			}, function(err) {
				reject(err);
			});
	});
}

/**
 * @promise LogoutPromise
 * @reject {ServerError}
 */

/**
 * Destroys the stored refresh token. It is safe to call this function when
 * already logged out.
 *
 * @return {LogoutPromise}
 */
function logout() {
	return new Promise(function(resolve, reject) {

		// If we aren't logged in, we can indicate success
		if (!isLoggedIn()) {
			resolve();
			return;
		}

		// If we don't have a refresh token, then it was already removed (perhaps
		// it was invalid), and we can proceed with the logout as normal.
		var refreshTokenId = window.localStorage.getItem('refreshTokenId');
		if (refreshTokenId === null) {
			window.sessionStorage.removeItem('loggedIn')
			resolve();
			return;
		}

		// We need to destroy the refresh token so it can't be used to generate
		// any new access tokens.
		destroyRefreshToken(refreshTokenId)
			.then(function() {
				window.localStorage.removeItem('refreshTokenId');
				window.sessionStorage.removeItem('loggedIn');
				resolve();
			}, function(err) {
				// A server error indicates that we're still technically logged
				// in; we'll have to try again later.
				if (err instanceof ServerError) {
					reject(err);
					return;
				}
				// Our failure wasn't due to a server error, so that means that
				// our refresh token was invalid. That's fine with us, we can
				// proceed with the logout as normal.
				//   err instanceof NotFoundError
				window.localStorage.removeItem('refreshTokenId');
				window.sessionStorage.removeItem('loggedIn');
				resolve();
			});
	});
}

/**
 * Checks if a user is currently logged in.
 *
 * @returns {boolean}
 */
function isLoggedIn() {
	return window.sessionStorage.getItem('loggedIn') === 'true';
}

/**
 * @promise AuthenticatedAjaxPromise
 * @reject {object, string, Error} jqXHR, textStatus, errorThrown
 * @reject {LoggedOutError}
 * @reject {ServerError}
 */

/**
 * Performs an authenticated $.ajax() call by attaching an "authorization" header
 * before making the request. The call is at most tried twice: once using the
 * cached access token, and if that fails, a new access token is requested and
 * the call is tried again.
 *
 * @param {string} url
 * @param {object} [settings]
 * @return {AuthenticatedAjaxPromise}
 */
function authenticatedAjax(url, settings) {
	return new Promise(function(resolve, reject) {
		function makeCall(accessToken) {
			settings = settings || {};
			settings.headers = settings.headers || {};
			settings.headers['authorization'] = 'Bearer ' + accessToken;
			return $.ajax(url, settings);
		}

		// First, get the access token.
		var p1 = getAccessToken();

		// If we don't have an access token, we can't make an authenticated ajax
		// call.
		p1.catch(reject);

		// Next, attempt to the make the ajax call.
		var p2 = p1.then(function(accessToken) {
			return makeCall(accessToken);
		});

		// If the call succeeds, then we're done.
		p2.then(resolve);

		// If the call doesn't succeed, it may be due to an expired/invalid
		// access token. We do try to catch these problems before making the
		// call, but the server gets to make the final decision.
		p2.catch(function(jqXHR) {

			// If the call didn't fail due to authentication, there's nothing we
			// can do about it.
			if (jqXHR.status !== 401) {
				reject.apply(this, arguments);
				return;
			}

			// It was an authentication problem, so let's request a new access
			// token.
			releaseAccessToken();
			var p3 = getAccessToken();

			// With our fresh access token, we'll try again. If there is still a
			// problem, there's nothing we can do about it.
			p3.then(makeCall)
				.then(resolve)
				.catch(reject);

			// If we still don't have an access token, we can't make an
			// authetnicated ajax. This may be due to a server error, or the user
			// may not be logged in anymore. In either case, getAccessToken()
			// would have handled the details, we just need to pass on the error.
			p3.catch(reject);
		});
	});
}

/**
 * @promise GetRefreshTokenPromise
 * @fulfill {string} refreshTokenId
 * @reject {InvalidCredentialsError} Only used if credentials are provided.
 * @reject {NotFoundError} Only used if credentials are not provided.
 * @reject {ServerError}
 */

/**
 * Gets or requests a refresh token. Storage is checked first, but if no token is
 * found, and credentials are provided, the credentials will be used to request a
 * new token.
 *
 * @param {string} [emailAddress]
 * @param {string} [password]
 * @return {GetRefreshTokenPromise}
 */
function getRefreshToken(emailAddress, password) {
	return new Promise(function(resolve, reject) {

		// Get the refresh token from cache
		var refreshTokenId = window.localStorage.getItem('refreshTokenId');
		if (refreshTokenId !== null) {
			resolve(refreshTokenId);
			return;
		}

		// If the refresh token wasn't in storage, and we weren't given any login
		// credentials, we can't get a refresh token.
		if (!emailAddress || !password) {
			reject(new NotFoundError());
			return;
		}

		// Create a new refresh token and cache it
		createRefreshToken(emailAddress, password)
			.then(function(refreshTokenId) {
				window.localStorage.setItem('refreshTokenId', refreshTokenId);
				resolve(refreshTokenId);
			}, reject);
	});
}

/**
 * @promise CreateRefreshTokenPromise
 * @fulfill {string} refreshTokenId
 * @reject {InvalidCredentialsError}
 * @reject {ServerError}
 */

/**
 * @param {string} emailAddress
 * @param {string} password
 * @return {CreateRefreshTokenPromise}
 */
function createRefreshToken(emailAddress, password) {
	return new Promise(function(resolve, reject) {
		var data = {
			emailAddress: emailAddress,
			password: password
		};
		$.post(makeApiUrl('/refresh-tokens'), JSON.stringify(data))
			.then(function(data, textStatus, jqXHR) {
				resolve(jqXHR.responseJSON.refreshTokenId);
			}, function(jqXHR) {
				if (jqXHR.status === 404) {
					reject(new InvalidCredentialsError());
					return;
				}
				reject(new ServerError());
			});
	});
}

/**
 * @promise DestroyRefreshTokenPromise
 * @reject {NotFoundError}
 * @reject {ServerError}
 */

/**
 * @param {string} refreshTokenId
 * @return {DestroyRefreshTokenPromise}
 */
function destroyRefreshToken(refreshTokenId) {
	return new Promise(function(resolve, reject) {
		$.ajax(makeApiUrl('/refresh-tokens/-'), {
			type: 'DELETE',
			headers: { 'authorization': 'Bearer ' + refreshTokenId }
		})
			.then(function() {
				resolve();
			}, function(jqXHR) {
				if (jqXHR.status === 401) {
					reject(new NotFoundError());
					return;
				}
				reject(new ServerError());
			});
	});
}

/**
 * @promise GetAccessTokenPromise
 * @fulfill {string} accessToken
 * @reject {LoggedOutError}
 * @reject {ServerError}
 */

/**
 * Gets a valid access token. If a cached token is available, that token will be
 * tried first. If the cached token is unavailable or invalid, a new token will
 * be generated using the refresh token, and will be cached for later (in
 * sessionStorage).
 *
 * @return {GetAccessTokenPromise}
 */
function getAccessToken() {
    return new Promise(function(resolve, reject) {
		var p1 = getRefreshToken();
		p1.then(function(refreshTokenId) {

			// Get the access token
			var accessToken = window.sessionStorage.getItem('accessToken');
			if (accessToken && isAccessTokenValid(accessTokenStr)) {
				resolve(accessToken);
				return;
			}

			// Try to generate a new access token
			var p2 = generateAccessToken(refreshTokenId)
			p2.then(function(accessToken) {
				// Save the new access token for later
				window.sessionStorage.setItem('accessToken', accessToken);
				resolve(accessToken);

			}, function(err) {
				// If there was a server error, then we should try again later.
				if (err instanceof ServerError) {
					reject(err);
					return;
				}

				// If there wasn't a server error, then that means our refresh
				// token is invalid. In other words:
				//   err instanceof InvalidCredentialsError
				logout()
					.then(function() {
						reject(new LoggedOutError());
					}, function() {
						// This should never happen
						throw new Error('logout failed after finding an invalid refresh token');
					});
			});
		}, function(err) {
			// getRefreshToken() was rejected
			if (err instanceof NotFoundError) {
				reject(new LoggedOutError());
				return;
			}
			// err instanceof ServerError
			reject(err);
		});
    });
}

/**
 * Releases the cached access token so that the next call to getAccessToken() is
 * forced to generate a fresh token.
 */
function releaseAccessToken() {
	window.sessionStorage.removeItem('accessToken');
}

/**
 * Checks whether an access token is a valid JWT. It also checks the expiration
 * value.
 *
 * @return {boolean}
 */
function isAccessTokenValid(accessToken) {

    // Decode the token
    var payload = null;
    try {
		payload = jwt_decode(accessToken);
    } catch(e) {
		return false;
    }

    // Check the expiration
    if (payload.exp) {
		var now = Math.floor(Date.now() / 1000);
		if (now >= payload.exp) {
			return false;
		}
    }

    return true;
}

/**
 * @promise GenerateAccessTokenPromise
 * @fulfill {string} accessToken
 * @reject {InvalidCredentialsError} Invalid refresh token.
 * @reject {ServerError}
 */

/**
 * @param {string} refreshTokenId
 * @return {GenerateAccessTokenPromise}
 */
function generateAccessToken(refreshTokenId) {
    return new Promise(function(resolve, reject) {
		$.ajax(makeApiUrl('/access-tokens'), {
			type: 'POST',
			headers: { 'authorization': 'Bearer ' + refreshTokenId },
		})
			.then(function(data, textStatus, jqXHR) {
				resolve(jqXHR.responseJSON.accessToken);
			}, function(jqXHR) {
				if (jqXHR.status === 401) {
					reject(new InvalidCredentialsError());
					return;
				}
				reject(new ServerError());
			});
    });
}
