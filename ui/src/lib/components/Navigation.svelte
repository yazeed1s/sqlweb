<script lang="ts">
	import {
		faCode,
		faDatabase
	} from '@fortawesome/free-solid-svg-icons';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import Fa from 'svelte-fa';

	let isHovered = false;

	onMount(() => {
		const sidebarIcon = document.querySelector('.sidebar-icon');

		sidebarIcon?.addEventListener('mouseenter', handleMouseEnter);
		sidebarIcon?.addEventListener('mouseleave', handleMouseLeave);

		return () => {
			sidebarIcon?.removeEventListener('mouseenter', handleMouseEnter);
			sidebarIcon?.removeEventListener('mouseleave', handleMouseLeave);
		};
	});

	const handleClick = (val: string): void => {
		switch (val) {
			case 'database':
				goto('/db');
				break;
			case 'settings':
				goto('/connection');
				break;
			case 'sql':
				goto('/editor');
				break;
		}
	};

	const handleMouseEnter = (): void => {
		isHovered = true;
	};

	const handleMouseLeave = (): void => {
		isHovered = false;
	};
</script>

<div>
	<div
		class="sidebar-icon"
		on:click={() => handleClick('database')}
		on:keydown={() => {}}
		on:keyup={() => {}}
		on:keypress={() => {}}
	>
		<span class="sidebar-tooltip">Database</span>
		<Fa icon={faDatabase} size="1.2rem" />
	</div>
	<!-- <div
		class="sidebar-icon"
		on:click={() => handleClick('settings')}
		on:keydown={() => {}}
		on:keyup={() => {}}
		on:keypress={() => {}}
	>
		<span class="sidebar-tooltip">Settings</span>
		<Fa icon={faCog} size="1.2rem" />
	</div> -->
	<div
		class="sidebar-icon"
		on:click={() => handleClick('sql')}
		on:keydown={() => {}}
		on:keyup={() => {}}
		on:keypress={() => {}}
	>
		<span class="sidebar-tooltip">SQL</span>
		<Fa icon={faCode} size="1.2rem" />
	</div>
</div>

<style>
	.sidebar-icon {
		/* position: relative; */
		cursor: pointer;
		outline: none;
		display: flex;
		align-items: center;
		justify-content: center;
		height: 2rem;
		width: 2rem;
		z-index: 3;
		margin: 0.9rem 0.3rem;
		background-color: rgb(var(--color-surface-800));
		color: rgb(var(--color-surface-200));
		cursor: pointer;
	}

	.sidebar-icon:hover {
		/* background-color: rgb(var(--color-surface-400)); */
		color: rgb(var(--color-primary-600));
		border-radius: 0.5rem;
		z-index: 3;
	}

	.sidebar-tooltip {
		position: absolute;
		width: auto;
		padding: 0.5rem;
		margin: 0.5rem;
		min-width: max-content;
		left: 2.5rem;
		border-radius: 0.375rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
		color: rgb(var(--color-surface-800));
		background-color: rgb(var(--color-tertiary-600));
		font-size: 0.75rem;
		font-weight: bold;
		transition: transform 0.1s linear;
		transform: scale(0);
		transform-origin: left;
		z-index: 2;
	}

	.sidebar-icon:hover .sidebar-tooltip {
		transform: scale(1);
		z-index: 1;
	}
</style>
