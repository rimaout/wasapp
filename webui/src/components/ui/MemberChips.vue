<script setup>
import UserAvatar from './UserAvatar.vue';

// Displays selected members as chips with an avatar, display name, and optional remove button.
defineProps({
    members:   { type: Array,   required: true }, // Array of member objects containing user details (id, name)
    removable: { type: Boolean, default: true  }, // Enables/disables the '×' remove button on each chip
    light:     { type: Boolean, default: false }, // Applies lighter color palette when true (useful for dark backgrounds)
});

// Declares custom events to emit back to parent component when actions occur
const emit = defineEmits(['remove']);
</script>

<template>
    <!-- Rendered only if member list is non-empty -->
    <div v-if="members.length > 0" class="member-chips">
        <!-- Loop through each member to render individual chip -->
        <span v-for="m in members" :key="m.id" class="member-chip" :class="{ light: light }">

            <!-- User avatar icon display -->
            <UserAvatar :userId="m.id" :displayName="m.name" :size="24" />

            <!-- Display member name with CSS truncation -->
            <span class="member-chip-name">{{ m.name }}</span>

            <!-- Remove button: passes member object to parent via 'remove' event -->
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
    background: var(--theme-bg-highlight);
    border-radius: 999px;
    padding: 3px 8px 3px 3px;
    color: var(--theme-fg);
}

.member-chip.light {
    background: var(--theme-overlay);
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
    color: var(--theme-fg-dark);
    line-height: 1;
    font-size: 1rem;
}

.member-chip-remove:hover {
    color: var(--theme-red);
}
</style>
