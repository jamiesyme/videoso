<template>
	<div class="login-view">
		<div class="container">
			<div class="short-container">
				<h2>Log In</h2>
				<form @submit.prevent="goToAdmin()">
					<input type="email" placeholder="Email address">
					<input type="password" placeholder="Password">
					<button type="submit">Log in</button>
				</form>

				<div class="note">
					<p><strong>Note:</strong></p>
					<p>This is a concept site; there is no way for a new user to register an account and log in.</p>
					<p>This page is simply available for administrative tasks, such as adding or removing content.</p>
					<p>If you'd like to contact the developer of this site, please visit <a href="https://jamiesyme.com">jamiesyme.com</a>.</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	export default {
		data () {
			return {};
		},

		methods: {
			goToAdmin () {
				this.$router.push('admin');
			},
		},
	}
</script>

<style lang="scss" scoped>
	.container {
		padding: 0;
	}

	.short-container {
		width: 33%;
		margin: 0 auto;
	}

	h2 {
		margin: 4rem 0 2rem;
	}

	button {
		display: block;
		width: 10rem;
		margin: 0 0 0 auto;
	}

	.note {
		//background-color: #ece6f2;
		//background-color: #DEF3FA;
		background-color: #E0F4E9;
		//border: 0.2rem solid #9b4dca;
		//border: 0.2rem solid #40bae3;
		border: 0.2rem solid #60C78F;
		border-radius: 0.4rem;
		padding: 1.6rem;

		p:last-of-type {
			margin-bottom: 1rem;
		}

		a {
			color: #60C78F;
			font-weight: bold;
		}
	}
</style>
