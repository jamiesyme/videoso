<template>
	<header class="view-header">
		<div :class="containerClasses">
			<div class="logo-wrapper">
				<router-link to="/">
					<img
						src="@/assets/logo.png"
						alt="Videoso logo">
				</router-link>
			</div>
			<div
				class="flex-grow"
				v-if="!isLargeContainer()">
			</div>
			<form
				class="search-wrapper"
				@submit.prevent="goToSearch()">
				<button
					type="button"
					class="condensed-search-button button button-clear"
					@click="activateSearch()">
					<IonSearch />
				</button>
				<div
					class="search-input-bg"
					v-show="searchActive">
				</div>
				<input
					type="text"
					ref="searchInput"
					placeholder="Search"
					:class="searchInputClasses"
					@blur="onSearchBlur()">
			</form>
			<div class="login-wrapper">
				<router-link
					to="/login"
					class="button button-clear"
					v-if="!loggedIn">
					Log in
				</router-link>
				<button
					class="button button-clear"
					@click="logout()"
					v-else>
					Log out
				</button>
			</div>
		</div>
	</header>
</template>

<script>
	import IonSearch from 'vue-ionicons/dist/md-search';
	import Auth from '@/auth';

	export default {
		components: {
			IonSearch,
		},

		data () {
			return {
				searchActive: false,
			};
		},

		methods: {
			activateSearch () {
				this.searchActive = true;
				this.$nextTick(() => {
					const elem = this.$refs.searchInput;
					if (elem) {
						elem.focus();
					}
				});
			},

			onSearchBlur () {
				this.searchActive = false;
			},

			isLargeContainer () {
				const bp = this.$breakpoint.name;
				return bp === 'large' || bp === 'xlarge';
			},

			goToSearch () {
				const q = this.$refs.searchInput.value;
				if (q) {
					this.$router.push({
						path: '/search',
						query: { q },
					});
					this.$refs.searchInput.blur();
				}
			},

			logout () {
				Auth.logout();
				if (this.$route.path.indexOf('/admin') === 0) {
					this.$router.push('/');
				}
			},
		},

		computed: {
			loggedIn () {
				return Auth.loggedIn();
			},

			containerClasses () {
				return {
					'container': true,
					'large-container': this.isLargeContainer(),
				};
			},

			searchInputClasses () {
				return {
					active: this.searchActive,
				};
			},
		},
	}
</script>

<style lang="scss" scoped>
	.view-header {
		height: 7.2rem;
		line-height: 7.2rem;
		border-bottom: 0.1rem solid #d6d6d6;
	}

	.container {
		display: flex;
	}

	.flex-grow {
		flex: 1;
	}

	.logo-wrapper {
		a {
			display: block;
			width: 14rem;
			height: 7.2rem;
			overflow: hidden;
		}

		img {
			width: auto;
			height: 14rem;
			position: relative;
			top: -3.4rem;
			left: -1.2rem;
		}
	}

	.search-wrapper {
		display: flex;
		align-items: center;
		margin: 0;

		input {
			display: none;
			margin: 0;
			background-color: #fafafa;
			z-index: 100;

			&:focus {
				background-color: white;
			}

			&.active {
				display: block;
				position: absolute;
				left: 1rem;
				width: calc(100% - 2rem);
			}
		}

		.search-input-bg {
			position: absolute;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0.1rem;
			background-color: white;
			z-index: 99;
		}

		.condensed-search-button {
			font-size: 2.0rem;
			padding: 0.3rem 1rem 0;
			margin: 0 0.5rem;
		}
	}

	.login-wrapper {
		padding: 0 0.4rem 0 0;
		font-weight: bold;
		text-align: right;

		.button {
			padding: 0 1rem;
		}
	}

	.large-container {
		.logo-wrapper {
			flex: 1 1 25%;
		}

		.search-wrapper {
			flex: 1 1 50%;

			input {
				display: block;

				&.active {
					position: static;
				}
			}

			.search-input-bg {
				display: none;
			}

			.condensed-search-button {
				display: none;
			}
		}

		.login-wrapper {
			flex: 1 1 25%;
		}
	}
</style>
