<script lang="ts">
	import go4 from '$lib/images/go4.png';
	import Fa from 'svelte-fa';
	import { faCog, faGear, faPowerOff } from '@fortawesome/free-solid-svg-icons';
	import { onMount } from 'svelte';
	import { fly, scale } from 'svelte/transition';
	import { connectedStore } from '$lib/store';
	import { goto } from '$app/navigation';
	import { quintOut } from 'svelte/easing';

	let show: boolean = false;
	let menu: any = null;
	let connected: boolean = false;
	let response: string = '';
	let schemas: string[] = [];

	connectedStore.subscribe((value) => {
		connected = value;
	});

	onMount(() => {
		const handleOutsideClick = (event: { target: any }) => {
			if (show && !menu.contains(event.target)) {
				show = false;
			}
		};
		const handleEscape = (event: { key: string }) => {
			if (show && event.key === 'Escape') {
				show = false;
			}
		};
		document.addEventListener('click', handleOutsideClick, false);
		document.addEventListener('keyup', handleEscape, false);
		return () => {
			document.removeEventListener('click', handleOutsideClick, false);
			document.removeEventListener('keyup', handleEscape, false);
		};
	});

	const handleClick = (val: number) => {
		switch (val) {
			case 0:
				disconnect();
				break;
			case 1:
				showSchemas();
				break;
		}
	};

	const routeToPage = (route: string): void => {
		goto(`/${route}`);
	};

	const fetchAndHandleResponse = async (url: string, method: string, body?: any): Promise<any> => {
		try {
			const options = {
				method,
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(body)
			};
			const res = await fetch(url, options);
			const json = await res.json();
			(res.ok) ? handleSuccessResponse('', json) : handleErrorResponse(json);
			return json;
		} catch (error) {
			console.error('An error occurred:', error);
		}
	};

	const disconnect = async (): Promise<void> => {
		await fetchAndHandleResponse('http://localhost:3000/disconnect', 'POST');
	};

	const showSchemas = async (): Promise<void> => {
		const json = await fetchAndHandleResponse('http://localhost:3000/schemas', 'GET');
		schemas = json.schemas;
		console.log(schemas);
	};

	const handleSuccessResponse = (route: string, json: any): void => {
		response = json.message;
		connectedStore.set(false);
		routeToPage(route);
	};

	const handleErrorResponse = (json: any): void => {
		response = json.error;
	};

	const handleKeyDown = (
		event: KeyboardEvent & { currentTarget: EventTarget & HTMLAnchorElement },
		index: number
	): void => {
		if (event.key === 'Enter') {
			handleClick(index);
		}
	};
</script>

<div class="top-navigation">
	<div class="dbit-container">
		<img src={go4} alt="Go" class="small-image" />
		<strong class="title">DBit</strong>
	</div>
	<div class="top-navigation-icon" bind:this={menu}>
		<button on:click={() => (show = !show)}>
			<Fa icon={faCog} />
		</button>
		{#if show}
			<div
				transition:fly={{ duration: 200, x: 100, easing: quintOut }}
				class="dropdown"
			>
				<!-- svelte-ignore a11y-invalid-attribute -->
				<a
					href="#"
					class={connected ? 'item' : 'disabled'}
					on:click={() => handleClick(0)}
					on:keydown={(event) => handleKeyDown(event, 0)}
				>
					<Fa icon={faPowerOff} style="margin-right: 0.5rem;" />
					<span>Disconnect</span>
				</a>
				<!-- svelte-ignore a11y-invalid-attribute -->
				<!-- <a
					href="#"
					class={connected ? 'item' : 'disabled'}
					on:click={() => handleClick(1)}
					on:keydown={(event) => handleKeyDown(event, 1)}
				>
					<Fa icon={faGear} style="margin-right: 0.5rem;" />
					<span>Change database</span>
				</a> -->
			</div>
		{/if}
	</div>
</div>

<style>
	.top-navigation {
		display: flex;
		flex-direction: row;
		align-items: center;
		justify-content: space-between;
		width: 100vw;
		height: 2.2rem;
		background-color: rgb(var(--color-surface-600));
		border-bottom: 1px solid rgb(var(--color-surface-500));
	}

	.top-navigation-icon {
		color: rgb(var(--color-text-primary));
		margin-right: 1.4rem;
		margin-left: 1rem;
		transition: color 0.3s ease-in-out;
		cursor: pointer;
	}

	.top-navigation-icon:hover {
		color: rgb(var(--color-primary-500));
	}

	.top-navigation-icon:first-child {
		margin-left: auto;
		margin-right: 2rem;
	}
	.title {
		text-align: center;
		position: relative;
		font-size: 16px;
		display: block;
		color: rgb(var(--color-surface-100));
	}
	.small-image {
		width: calc(1em + 0.9rem);
		height: auto;
		margin-left: 10px;
	}
	.dbit-container {
		display: flex;
		align-items: center;
		position: relative;
	}

	.top-navigation:before {
		content: '';
		position: absolute;
		left: 0;
		bottom: 0;
		width: 100%;
		height: 1px;
	}

	.dropdown {
		transform-origin: top right;
		position: absolute;
		right: 0px;
		width: 12rem;
		padding-top: 0.5rem;
		padding-bottom: 0.5rem;
		margin-top: 0.5rem;
		margin-right: 0.85rem;
		background-color: rgb(var(--color-surface-500));
		border-radius: 0.25rem;
		box-shadow: 0px 2px 10px rgba(0, 0, 0, 0.452);
		z-index: 5;
	}

	.item {
		display: flex;
		align-items: center;
		padding: 0.5rem 1rem;
		color: rgb(var(--color-surface-100));
	}

	.item span:hover {
		color: rgb(var(--color-primary-500));
	}

	.disabled {
		display: flex;
		align-items: center;
		padding: 0.5rem 1rem;
		pointer-events: none;
	}

	.disabled span:hover {
		color: #747474;
		pointer-events: none;
	}
</style>
