<script lang="ts">
	import 'iconify-icon';
	import type { Table } from '../store';
	import { connectedStore, selectedTable, tableDataStore } from '../store';
	import { afterUpdate, onDestroy, onMount } from 'svelte';
	import { scale, slide } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';

	let tableData: { tables: Table[] } = { tables: [] };
	let connected: boolean = false;
	let resizeHandle: HTMLElement | null = null;
	let resizing: boolean = false;
	let initialWidth: number = 150;
	let expandedTable: string | null = null;
	let expandedColumn: string | null = null;

	connectedStore.subscribe((value) => {
		connected = value;
	});

	const handleClick = (table: string) => {
		selectedTable.set({
			tableName: table
		});
	};

	const handleToggleColumns = (tableName: string | null) => {
		expandedTable = expandedTable === tableName ? null : tableName;
	};

	const handleToggleColumnsInfo = (colName: string | null) => {
		expandedColumn = expandedColumn === colName ? null : colName;
	};

	onMount(() => {
		return tableDataStore.subscribe((data) => {
			tableData = data;
		});
	});

	afterUpdate(() => {
		if (resizeHandle) {
			resizeHandle.addEventListener('mousedown', startResize);
			window.addEventListener('mousemove', handleResize);
			window.addEventListener('mouseup', stopResize);
		}
	});

	onDestroy(() => {
		if (resizeHandle) {
			resizeHandle.removeEventListener('mousedown', startResize);
			window.removeEventListener('mousemove', handleResize);
			window.removeEventListener('mouseup', stopResize);
		}
	});

	const startResize = (event: MouseEvent): void => {
		resizing = true;
		initialWidth = event.clientX;
	};

	const handleResize = (event: MouseEvent): void => {
		if (!resizing) return;
		const listNav = document.querySelector('.list-nav') as HTMLElement;
		const movementX = event.clientX - initialWidth;
		const newWidth = parseInt(getComputedStyle(listNav).width || '0') + movementX;
		listNav.style.width = `${newWidth}px`;
		initialWidth = event.clientX;
	};

	const stopResize = (): void => {
		resizing = false;
	};
</script>

{#if connected}
	<nav class="list-nav" bind:this={resizeHandle}>
		<ul class="table-names">
			{#each tableData.tables as table}
				<li transition:slide={{ duration: 100, axis: 'y' }}>
					<div class="table-item">
						<a href="/db/tables" on:click={() => handleClick(table.name)}>
							<span class="tbl-icon-container">
								<iconify-icon icon="mdi:table" />
							</span>
							<span class="table-name">{table.name}</span>
						</a>
						<div class="toggle-button-container">
							<button
								class="toggle-button"
								on:click={() => handleToggleColumns(table.name)}
							>
								<iconify-icon icon="mdi:chevron-down" />
							</button>
						</div>
					</div>
					{#if expandedTable === table.name}
						<ul
							class="cols-names"
							transition:slide={{
								delay: 50,
								duration: 500,
								easing: quintOut,
								axis: 'y'
							}}
						>
							{#each table.columns as column}
								<li>
									<span class="cols-icon-container">
										<iconify-icon icon="mdi:table-column" />
									</span>
									<!-- svelte-ignore a11y-click-events-have-key-events -->
									<span on:click={() => handleToggleColumnsInfo(column.field)}>
										{column.field}
									</span>
									{#if expandedColumn === column.field}
										<div
											class="cols-info"
											transition:slide={{
												delay: 50,
												duration: 500,
												easing: quintOut,
												axis: 'y'
											}}
										>
											{#if column.type}
												<span>Type: {column.type}</span>
											{/if}
											{#if column.key}
												<span>Key: {column.key}</span>
											{/if}
											{#if column.refrenced_table}
												<span
													>Referenced table: {column.refrenced_table}</span
												>
											{/if}
											{#if column.refrenced_column}
												<span
													>Referenced column: {column.refrenced_column}</span
												>
											{/if}
											{#if column.constraint_name}
												<span
													>Constraint name: {column.constraint_name}</span
												>
											{/if}
										</div>
									{/if}
								</li>
							{/each}
						</ul>
					{/if}
				</li>
			{/each}
		</ul>
	</nav>
{/if}

<style>
	.list-nav {
		margin-top: 0.75vh;
		width: 150px;
		max-width: 220px;
		min-width: 10px;
		list-style: none;
		height: 100vh;
		transition: left 0.4s ease;
	}

	.table-names {
		z-index: 0;
		position: absolute;
		padding-left: 0.5vw;
		overflow-y: auto;
		overflow-x: auto;
		height: calc(100vh - 50px);
	}

	.table-item {
		display: flex;
		align-items: center;
	}

	.table-item a {
		display: flex;
		align-items: center;
		text-decoration: none;
		color: rgb(var(--color-surface-100));
		padding: 5px;
	}

	.table-item a span:hover {
		color: rgb(var(--color-primary-500));
	}

	.toggle-button-container {
		margin-left: 5px;
	}

	.tbl-icon-container {
		display: inline-flex;
		align-items: center;
		margin-right: 5px;
	}

	.cols-icon-container {
		display: inline-flex;
		align-items: center;
		margin-right: 2px;
	}

	.cols-names {
		overflow: hidden;
		text-align: left;
		margin-left: 1.2rem;
	}

	.cols-names li span {
		color: rgb(var(--color-text-secondary));
		cursor: pointer;
	}

	.cols-names li span:hover {
		color: rgb(var(--color-primary-500));
	}

	.cols-info {
		overflow: hidden;
		text-align: left;
		margin-left: 1.9rem;
		display: flex;
		flex-direction: column;
		margin-top: 5px;
	}
</style>
