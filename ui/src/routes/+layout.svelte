<script lang="ts">
	import '../app.css';
	import '../app.postcss';
	import { page } from '$app/stores';
	import Navigation from '$lib/components/Navigation.svelte';
	import PageHeader from '$lib/components/PageHeader.svelte';
	import { connectedStore } from '$lib/store';
	import Connection from '$lib/components/Connection.svelte';
	import Error from './+error.svelte';
	import { onDestroy } from 'svelte';
	let connected = false;
	const unsubscribe = connectedStore.subscribe((value) => {
		connected = value;
	});
	onDestroy(() => {
		unsubscribe();
	});
</script>

{#if !$page.error}
	<body>
		<header>
			<PageHeader />
		</header>
		<div class="content">
			<aside class="left">
				{#if connected}
					<Navigation />
				{/if}
			</aside>
			<main>
				{#if !connected}
					<Connection />
				{/if}
				<slot />
			</main>
		</div>
	</body>
{:else}
	<Error />
{/if}

<style>
	* {
		margin: 0;
		box-sizing: border-box;
	}
	body {
		background-color: rgb(var(--color-surface-800));
		height: 100%;
		width: 100%;
		text-align: center;
		display: flex;
		flex-direction: column;
	}
	header {
		background-color: rgb(var(--color-surface-700));
		padding: 0 0 0 0;
	}

	.left {
		padding: 0 0 3em 0;
		flex: 0 1 1px;
		border-right: 1px solid rgb(var(--color-surface-500));
	}
	.content {
		display: flex;
		flex: 0;
		width: calc(100vw - 2vw);
	}
	main {
		padding: 0 0 0 0;
		flex: 5 5 150px;
		width: 100%;
	}

	@media all and (max-width: 550px) {
		.content {
			flex-direction: column;
		}
		main {
			padding: 5em 0 5em 0;
		}
	}
</style>
