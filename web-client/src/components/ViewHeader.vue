<template>
	<header class="view-header row">
		<div class="logo-wrapper column column-25">
			<img alt="Videoso logo" src="@/assets/logo.png">
		</div>
		<div class="search-wrapper column">
			<input
				type="text"
				placeholder="Search..."
				v-if="showSearch">
		</div>
		<div class="login-wrapper column column-25">
			<router-link to="#">Log in</router-link>
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
		height: 10rem;
		line-height: 10rem;
	}

	.logo-wrapper {
		img {
			width: auto;
			height: 100%;
		}
	}

	.login-wrapper {
		font-weight: bold;
		text-align: right;
	}
</style>
