<template>
	<div class="label-input">
		<input
			type="text"
			ref="input"
			:value="value"
			@keydown="checkEnter($event)"
			@input="$emit('input', $event.target.value)"
			@blur="stopEditing()"
			v-if="editing">
		<div
			class="label"
			@click="startEditing()"
			v-else>
			{{ value }}
		</div>
	</div>
</template>

<script>
	export default {
		props: {
			value: String,
		},

		data () {
			return {
				editing: false,
			};
		},

		methods: {
			startEditing () {
				if (this.editing) {
					return;
				}
				this.editing = true;
				this.$nextTick(() => {
					this.$refs.input.focus();
					this.$refs.input.select();
				});
			},

			stopEditing () {
				if (!this.editing) {
					return;
				}
				this.editing = false;
			},

			checkEnter (ev) {
				if (ev.key === 'Enter') {
					this.stopEditing();
				}
			},
		},
	}
</script>

<style lang="scss" scoped>
	input {
		margin: 0;
	}

	.label:hover {
		cursor: text;
	}
</style>
