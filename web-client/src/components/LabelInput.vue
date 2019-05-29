<template>
	<div class="label-input">
		<input
			ref="input"
			:list="list"
			:type="type"
			:value="value"
			@keydown="checkEnter($event)"
			@input="$emit('input', $event.target.value)"
			@blur="stopEditing()"
			v-if="editing">
		<div
			class="label"
			@click="startEditing()"
			v-else>
			<div class="label-value" v-if="value">{{ value }}</div>
			<slot
				name="empty"
				v-else>
				<div class="empty-label"></div>
			</slot>
		</div>
	</div>
</template>

<script>
	export default {
		props: {
			list: String,
			type: {
				type: String,
				default: 'text',
			},
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
	.label-value {
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.empty-label {
		min-height: 2.4rem;
	}
</style>
