<script lang="ts">
	import { page } from '$app/stores';
	import DbNavigation from '$lib/components/DbNavigation.svelte';
	import { connectedStore } from '$lib/store';
	import Error from '../+error.svelte';
	let connected = false;
	connectedStore.subscribe((value) => {
		connected = value;
	});
</script>

{#if !$page.error}
	{#if connected}
		<body>
			<div class="content">
				<aside class="left">
					<DbNavigation />
				</aside>
				<main>
					<slot />
				</main>
			</div>
		</body>
	{/if}
{:else}
	<Error />
{/if}

<style>
	* {
		margin: 0;
		box-sizing: border-box;
	}
	body {
		height: 100%;
		/* width: 97vw; */
		width: 100%;
		background-color: rgb(var(--color-surface-800));
		text-align: center;
		display: flex;
		flex-direction: column;
	}

	.left {
		padding: 0 0 0 0;
		flex: 0 0 0.5px;
		border-right: 1px solid rgb(var(--color-surface-500));
		flex-shrink: 0;
		background-color: rgb(var(--color-surface-800));
	}

	.content {
		display: flex;
		flex: 1;
	}

	main {
		display: flex;
		flex: 1;
		padding: 0.5em 0 1em 0;
		/* flex: 5 5 150px; */
		background-color: rgb(var(--color-surface-800));

		/* max-width: calc(100%); */
		width: 100%;
		min-width: calc(100% - 220px);
		z-index: 2;
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
