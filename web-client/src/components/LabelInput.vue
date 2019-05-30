<template>
	<div class="label-input">
		<div
			ref="inputWrapper"
			class="input-wrapper"
			@blur.capture="stopEditing()"
			v-if="editing">
			<slot name="input">
				<input
					type="datetime-local"
					:value="dateToDateTimeLocalStr(value)"
					@keydown="checkEnter($event)"
					@input="$emit('input', dateTimeLocalStrToDate($event.target.value))"
					v-if="type === 'datetime-local'">
				<input
					:list="list"
					:type="type"
					:value="value"
					@keydown="checkEnter($event)"
					@input="$emit('input', $event.target.value)"
					v-else>
			</slot>
		</div>
		<div
			class="label-wrapper"
			@click="startEditing()"
			v-else>
			<slot name="label">
				<div
					class="label"
					v-if="!isValueEmpty(value)">
					{{ value }}
				</div>
				<slot
					name="empty"
					v-else>
					<div class="label empty-label"></div>
				</slot>
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
			value: [String, Date, Number],
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
					const elem = this.$refs.inputWrapper.querySelector('input');
					if (elem) {
						elem.focus();
						elem.select();
					}
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

			isValueEmpty (value) {
				return (value === null ||
				        value === '' ||
				        typeof value === 'undefined');
			},

			dateToDateTimeLocalStr (date) {
				if (date instanceof Date) {
					date = new Date(date);
					const offset = date.getTimezoneOffset();
					const msPerMin = 60 * 1000;
					date.setTime(date.getTime() - offset * msPerMin);
					return date.toISOString().substr(0, 16);
				} else {
					return null;
				}
			},

			dateTimeLocalStrToDate (str) {
				const ms = Date.parse(str);
				if (typeof ms === 'number') {
					return new Date(ms);
				} else {
					return null;
				}
			},
		},
	}
</script>

<style lang="scss" scoped>
	input {
		margin: 0;
	}

	.label-wrapper:hover {
		cursor: text;
	}

	.label {
		overflow: hidden;
		text-overflow: ellipsis;
		&.empty-label {
			min-height: 2.4rem;
		}
	}
</style>
