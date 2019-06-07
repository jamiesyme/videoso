<template>
	<header class="view-header">
		<div class="container">
			<div class="logo-wrapper">
				<router-link to="/">
					<img
						src="@/assets/logo.png"
						alt="Videoso logo">
				</router-link>
			</div>
			<form
				class="search-wrapper"
				@submit.prevent="goToSearch()">
				<input
					type="text"
					ref="searchInput"
					placeholder="Search"
					v-show="showSearch">
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
	import Auth from '@/auth';

	export default {
		props: {
			showSearch: {
				type: Boolean,
				default: true,
			},
		},

		methods: {
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

	.logo-wrapper {
		flex: 1 1 25%;

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
		flex: 1 1 50%;

		input {
			margin: 0;
			background-color: #fafafa;

			&:focus {
				background-color: white;
			}
		}
	}

	.login-wrapper {
		flex: 1 1 25%;
		padding: 0 0.4rem 0 0;
		font-weight: bold;
		text-align: right;

		.button {
			padding: 0 1rem;
		}
	}
</style>
