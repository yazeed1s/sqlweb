<script lang="ts">
	// import Editor from '$lib/components/Editor.svelte';
	import { page } from '$app/stores';
	import Error from '../+error.svelte';
</script>

{#if !$page.error}
	<body>
		<div class="content">
			<main>
				<slot />
			</main>
		</div>
	</body>
{:else}
	<Error />
{/if}

<style>
	* {
		margin: 0px;
		box-sizing: border-box;
	}
	body {
		height: 100%;
		width: 100%;
		text-align: center;
		display: flex;
		flex-direction: column;
	}

	.content {
		width: 100%;
		display: flex;
		flex: 1;
	}
	main {
		width: 100%;
		padding: 1em 0 0em 0;
		flex: 5 5 150px;
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
