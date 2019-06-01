<template>
	<header class="view-header">
		<div class="container">
			<div class="logo-wrapper">
				<router-link to="/">
					<img alt="Videoso logo" src="@/assets/logo.png">
				</router-link>
			</div>
			<div class="search-wrapper">
				<input
					type="text"
					placeholder="Search"
					v-if="showSearch">
			</div>
			<div class="login-wrapper">
				<router-link to="/login" class="button button-clear">
					Log in
				</router-link>
			</div>
		</div>
	</header>
</template>

<script>
	export default {
		props: {
			showSearch: {
				type: Boolean,
				default: true,
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

	.logo-wrapper {
		position: absolute;
		left: 0;
		width: 25%;
		padding: 0 0 0 0.4rem;

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
			left: -1rem;
		}
	}

	.search-wrapper {
		position: absolute;
		left: 25%;
		width: 50%;

		input {
			margin: 0;
			background-color: #fafafa;

			&:focus {
				background-color: white;
			}
		}
	}

	.login-wrapper {
		position: absolute;
		right: 0;
		width: 25%;
		padding: 0 0.4rem 0 0;
		font-weight: bold;
		text-align: right;

		.button {
			padding: 0 1rem;
		}
	}
</style>
