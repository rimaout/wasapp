// Selected members shown as wrapping chips with avatar, name and × remove.
<script setup>
import UserAvatar from './UserAvatar.vue';

defineProps({
	members: { type: Array, required: true },
	removable: { type: Boolean, default: true },
	light: { type: Boolean, default: false },
});
const emit = defineEmits(['remove']);
</script>

<template>
	<div v-if="members.length > 0" class="member-chips">
		<span v-for="m in members" :key="m.id" class="member-chip" :class="{ light: light }">
			<UserAvatar :userId="m.id" :displayName="m.name" :size="24" />
			<span class="member-chip-name">{{ m.name }}</span>
			<button v-if="removable" type="button" class="member-chip-remove" @click="emit('remove', m)">×</button>
		</span>
	</div>
</template>

<style scoped>
.member-chips {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
}

.member-chip {
	display: inline-flex;
	align-items: center;
	gap: 6px;
	background: var(--tn-bg-highlight);
	border-radius: 999px;
	padding: 3px 8px 3px 3px;
	color: var(--tn-fg);
}

.member-chip.light {
	background: var(--tn-overlay);
}

.member-chip-name {
	font-size: 0.85rem;
	max-width: 120px;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.member-chip-remove {
	background: none;
	border: none;
	padding: 0;
	cursor: pointer;
	color: var(--tn-fg-dark);
	line-height: 1;
	font-size: 1rem;
}

.member-chip-remove:hover {
	color: var(--tn-red);
}
</style>
