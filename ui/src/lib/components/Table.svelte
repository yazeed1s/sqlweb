<script lang="ts">
	import { connectedStore, schemaName, selectedTable } from '../store';
	import { onMount } from 'svelte';
	import { flip } from 'svelte/animate';
	import { fade, fly, scale } from 'svelte/transition';
	import {
		faCheck,
		faDatabase,
		faEllipsisVertical,
		faFileCsv,
		faFileExport,
		faRotateLeft,
		faX
	} from '@fortawesome/free-solid-svg-icons';
	import Fa from 'svelte-fa';
	import { quintOut } from 'svelte/easing';

	const rowsPerPage: number = 300;
	let tableElement: HTMLTableElement;
	let show: boolean = false;
	let menu: any = null;
	let header: string[] = [];
	let rows: Array<any> = [];
	let cols: Array<ResponseCols> = [];
	let isLoading: boolean = true;
	let edit: boolean = false;
	let showIcons: boolean = false;
	let connected: boolean = false;
	let currentPage: number = 1;
	let totalRows: number;
	let totalPages: number;
	let editingRowIndex: number = -1;
	let editingColIndex: number = -1;
	let originalValue: string = '';
	let schema: string = '';
	let n_rows: string;
	let n_cols: string;
	let mb_size: string;
	let tableResponse: any;
	let req = { query: '' };
	let selectTableName: { selectedTable: string } = { selectedTable: '' };

	interface ResponseCols {
		field: string;
		type: string;
		key: string;
	}

	interface PrimaryKey {
		cellValue: string | null;
		headerValue: string | null;
	}

	interface RowValues {
		cellValue: string | null | undefined;
		headerValue: string | null | undefined;
		parentColumn: string | null;
		editedCellValue: string | null;
		tableName: string | null;
	}

	connectedStore.subscribe((value) => (connected = value));
	schemaName.subscribe((value) => (schema = value.name));

	onMount(() => {
		const unsubscribeSelectedTable = selectedTable.subscribe((data) => {
			const tableName = data.tableName;
			selectTableName.selectedTable = data.tableName;
			console.log(tableName);
		});
		const handleOutsideClick = (event: { target: any }) => {
			if (show && !menu.contains(event.target)) {
				show = false;
			}
		};
		document.addEventListener('click', handleOutsideClick, false);
		return () => {
			unsubscribeSelectedTable();
			document.addEventListener('click', handleOutsideClick, false);
		};
	});

	const editRowRequest = async (values: RowValues): Promise<void> => {
		const res = await fetch('http://localhost:3000/update', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(values)
		});
		const json = await res.json();
		console.log(json);
	};

	const exportRequest = async (choice: string): Promise<void> => {
		let url =
			choice === 'sql'
				? `http://localhost:3000/${choice}/export`
				: `http://localhost:3000/${choice}/export?name=${encodeURIComponent(
						selectTableName.selectedTable
				  )}`;

		let fileName =
			choice === 'sql'
				? `${schema}.sql`
				: choice === 'csv'
				? `${selectTableName.selectedTable}.csv`
				: `${selectTableName.selectedTable}.json`;
		const res = await fetch(url, { method: 'GET' });
		if (!res.ok) {
			console.error('Export request failed.');
			return;
		}
		console.log(res.headers);
		// const filename = res.headers.get('Filename');
		const blob = await res.blob();
		downloadFile(blob, fileName);
	};

	const downloadFile = (blob: Blob, file: string | null): void => {
		// Create a URL for the blob
		const uurl = window.URL.createObjectURL(blob);
		// Create a temporary anchor element to trigger the download
		const a = document.createElement('a');
		a.href = uurl;
		console.log(file);
		if (file) a.download = file;
		a.style.display = 'none';
		document.body.appendChild(a);
		// Click the anchor element to trigger the download
		a.click();
		// Clean up
		window.URL.revokeObjectURL(uurl);
		document.body.removeChild(a);
	};

	const initializeTableData = async (): Promise<void> => {
		header = [];
		rows = [];
		cols = [];
		let url = `http://localhost:3000/table?name=${encodeURIComponent(
			selectTableName.selectedTable
		)}&page=${currentPage}&perPage=${rowsPerPage}`;
		try {
			const response = await fetch(url, {
				method: 'GET',
				headers: { 'Content-Type': 'application/json' }
			});
			tableResponse = await response.json();
			console.log(tableResponse.data.table.data);
			if (tableResponse.data.table.data !== null) {
				header = Object.keys(tableResponse.data.table.data[0]);
				rows = tableResponse.data.table.data.map((row: any) => Object.values(row));
			} else {
				header = tableResponse.data.table.columns.map((column: any) => column.field);
			}
			cols = tableResponse.data.table.columns;
			n_rows = tableResponse.data.table.n_rows;
			n_cols = tableResponse.data.table.n_columns;
			mb_size = tableResponse.data.table.size_mb;
			totalRows = tableResponse.data.total_rows;
			totalPages = tableResponse.data.total_pages;
			console.log('header:', header);
			console.log('rows', rows);
			console.log('rows[0]', rows.at(0));
			console.log('columns:', cols);
			console.log('total_rows:', totalRows);
			console.log('total_pages:', totalPages);
		} catch (error) {
			console.error('Error fetching table data:', error);
		}
	};

	$: {
		if (selectTableName.selectedTable) {
			isLoading = true;
			currentPage = 1;
			initializeTableData();
			setTimeout(() => {
				isLoading = false;
			}, 1000);
		}
	}

	const reloadTable = (): void => {
		isLoading = true;
		initializeTableData();
		setTimeout(() => {
			isLoading = false;
		}, 1000);
	};

	const getMatchingCols = (header: string): { key: string; type: string }[] => {
		return cols
			.filter((col) => col.field === header)
			.map((col) => ({ key: col.key, type: col.type }));
	};

	const handleKeyDown = (
		event: KeyboardEvent & { currentTarget: EventTarget & HTMLDivElement }
	): void => {
		if (event.key === 'Enter') {
			reloadTable();
		}
	};

	const nextPage = (): void => {
		if (currentPage < totalPages) {
			currentPage++;
			reloadData();
		}
	};

	const prevPage = (): void => {
		if (currentPage > 1) {
			currentPage--;
			reloadData();
		}
	};

	const goToPage = (page: number): void => {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
			reloadData();
		}
	};

	const reloadData = async (): Promise<void> => {
		await initializeTableData();
	};

	const enableEditing = (rowIndex: number, colIndex: number): void => {
		if (edit == true) {
			editingRowIndex = rowIndex;
			editingColIndex = colIndex;
			rows[rowIndex].editing = true;
			originalValue = rows[rowIndex][colIndex];
		}
	};

	const cancelEditing = (): void => {
		rows[editingRowIndex][editingColIndex] = originalValue;
		rows[editingRowIndex].editing = false;
		editingRowIndex = -1;
		editingColIndex = -1;
		switchButton();
	};

	const findPrimaryKey = (rowIndex: number): PrimaryKey | undefined => {
		let arr: Array<string> = rows[rowIndex];
		for (let i = 0; i < arr.length; i++) {
			for (let j = 0; j < cols.length; j++) {
				if (header[i] === cols.at(j)?.field && cols.at(j)?.key === 'PRI') {
					const primaryKey: PrimaryKey = {
						cellValue: arr[i],
						headerValue: header[i]
					};
					console.log('obj col is: ', primaryKey.headerValue);
					console.log('obj primary key is: ', primaryKey.cellValue);
					return primaryKey;
				}
			}
		}
		return undefined;
	};

	const editRow = async (): Promise<void> => {
		switchButton();
		const primaryRecord = findPrimaryKey(editingRowIndex);
		const values: RowValues = {
			// for some weird reasons, cellValue gets casted to number
			cellValue: primaryRecord?.cellValue?.toString(),
			editedCellValue: rows[editingRowIndex][editingColIndex],
			headerValue: primaryRecord?.headerValue,
			parentColumn: header[editingColIndex],
			tableName: selectTableName.selectedTable
		};
		console.log(JSON.stringify(values));
		await editRowRequest(values);
	};

	const switchButton = (): void => {
		edit = !edit;
		showIcons = !showIcons;
		let button = document.getElementById('edit-button');
		edit
			? // eslint-disable-next-line svelte/valid-compile
			  (button!.style.cssText = 'background-color:  rgb(var(--color-success-700));')
			: (button!.style.cssText = 'background-color:  rgb(var(--color-surface-800));');
	};

	const handleExport = (choice: number): void => {
		switch (choice) {
			case 0:
				exportRequest('json');
				break;
			case 1:
				exportRequest('sql');
				break;
			case 2:
				exportRequest('csv');
				break;
		}
	};
</script>

{#if connected}
	<body>
		<header>
			<div class="header">
				{#if selectTableName.selectedTable}
					<div class="left-section">
						<h2 class="path">{schema} / {selectTableName.selectedTable}</h2>
						<div
							style="cursor: pointer;"
							on:click={() => reloadTable()}
							on:keydown={(event) => handleKeyDown(event)}
						>
							<Fa icon={faRotateLeft} />
						</div>
						<button style="" id="edit-button" on:click={() => switchButton()}>
							Edit</button
						>
						{#if showIcons}
							<div
								style="cursor: pointer; margin-left: 1rem;"
								on:click={() => editRow()}
								on:keydown={(event) => handleKeyDown(event)}
							>
								<Fa icon={faCheck} color="rgb(var(--color-success-500))" />
							</div>
							<div
								style="cursor: pointer; margin-left: 1rem;"
								on:click={() => cancelEditing()}
								on:keydown={(event) => handleKeyDown(event)}
							>
								<Fa icon={faX} color="rgb(var(--color-error-500))" />
							</div>
						{/if}
					</div>
				{/if}
				<div class="info">
					<h2>Rows: {totalRows}</h2>
					<h2>Columns: {n_cols}</h2>
					<h2>Size: {mb_size}</h2>
				</div>
			</div>
		</header>
		{#if isLoading}
			<div class="loader" />
		{:else}
			<div class="container">
				<div class="table-container" transition:fade={{ duration: 100 }}>
					<table bind:this={tableElement}>
						<thead class="sticky-header">
							<tr>
								{#each header as h}
									<th>
										{h}
									</th>
								{/each}
							</tr>
						</thead>
						<tbody>
							{#each rows as row, rowIndex}
								<tr>
									{#each row as value, colIndex}
										<!-- svelte-ignore a11y-click-events-have-key-events -->
										<td on:click={() => enableEditing(rowIndex, colIndex)}>
											{#if row.editing && rowIndex === editingRowIndex && colIndex === editingColIndex && row[colIndex] !== null && edit}
												<input
													type="text"
													bind:value={row[colIndex]}
													style="width: calc(1em * {row[colIndex].length +
														1});"
												/>
											{:else}
												{row[colIndex]}
											{/if}
										</td>
									{/each}
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/if}

		<footer>
			<div class="pagination">
				<button on:click={prevPage} disabled={currentPage === 1}>&larr;</button>
				<div class="page-number">{currentPage} / {totalPages}</div>
				<button on:click={() => nextPage()} disabled={currentPage === totalPages}
					>&rarr;</button
				>
			</div>
			<div class="_export" bind:this={menu}>
				<button on:click={() => (show = !show)}>
					<Fa icon={faEllipsisVertical} />
				</button>
				{#if show}
					<div
						transition:fly={{ duration: 200, y: 100, easing: quintOut }}
						class="dropdown"
					>
						<!-- svelte-ignore a11y-invalid-attribute -->
						<a
							href="#"
							class={connected ? 'item' : 'disabled'}
							on:click={() => {
								handleExport(0);
							}}
						>
							<Fa icon={faFileExport} style="margin-right: 0.5rem;" />
							<span>Export current table to JSON</span>
						</a>
						<!-- svelte-ignore a11y-invalid-attribute -->
						<a
							href="#"
							class={connected ? 'item' : 'disabled'}
							on:click={() => {
								handleExport(1);
							}}
						>
							<Fa icon={faDatabase} style="margin-right: 0.5rem;" />
							<span>Export db schema to SQL</span>
						</a>
						<!-- svelte-ignore a11y-invalid-attribute -->
						<a
							href="#"
							class={connected ? 'item' : 'disabled'}
							on:click={() => {
								handleExport(2);
							}}
						>
							<Fa icon={faFileCsv} style="margin-right: 0.5rem;" />
							<span>Export current table to CSV</span>
						</a>
					</div>
				{/if}
			</div>
		</footer>
	</body>
{/if}

<style>
	* {
		box-sizing: border-box;
	}

	body {
		height: 100%;
		width: 95.6vw;
		max-width: 100%;
		display: flex;
		flex-direction: column;
		text-align: center;
	}

	header {
		width: 100%;
	}

	.container {
		min-width: 100%;
		/* height: 93%; */
	}

	footer {
		background-color: rgb(var(--color-surface-800));
		position: fixed;
		left: 0;
		bottom: 0;
		width: 100%;
		height: 2em;
	}

	footer .pagination {
		flex: 0 1 auto;
		position: absolute;
		left: 50%;

		display: flex;
		align-items: center;
		align-content: center;
	}

	footer ._export {
		flex: 0 1 auto;
		position: absolute;
		margin-left: auto;
		right: 1%;
	}

	.dropdown {
		transform-origin: top;
		position: absolute;
		bottom: 2.4rem;
		right: 0;
		width: 17rem;
		background-color: rgb(var(--color-surface-500));
		border-radius: 6px;
		box-shadow: 0px 2px 10px rgba(0, 0, 0, 0.452);
		z-index: 2;
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

	tbody input[type='text'] {
		background-color: rgb(var(--color-surface-600));
		padding: 4px 8px;
		border-radius: 2px;
		border: none;
		box-shadow: none;
	}

	button {
		border-radius: 4px;
		padding: 5px 10px;
		background-color: rgb(var(--color-surface-800));
		color: rgb(var(--color-surface-100));
		cursor: pointer;
		transition: background-color 0.3s;
	}

	.pagination button {
		border-radius: 4px;
		height: 2.2rem;
		background-color: rgb(var(--color-surface-700));
		color: rgb(var(--color-surface-100));
		cursor: pointer;
		transition: background-color 0.3s;
	}

	#edit-button {
		border-radius: 7px;
		margin-left: 1em;
	}

	#edit-button:hover {
		background-color: rgb(var(--color-success-600));
	}

	button:disabled {
		background-color: rgb(var(--color-surface-700));
		color: rgb(var(--color-surface-300));
		cursor: default;
	}

	button:hover:not(:disabled) {
		background-color: rgb(var(--color-surface-600));
	}

	.page-number {
		margin: 0 20px;
		font-size: 14px;
		color: rgb(var(--color-text-secondary));
	}

	.header {
		/* padding: 0 0 0 0; */
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 100%;
		max-width: 100%;
		color: rgb(var(--color-text-primary));
	}

	.left-section {
		display: flex;
		align-items: center;
		background-color: rgb(var(--color-surface-800));
	}

	.table-container {
		margin-top: 0.3em;
		margin-bottom: 1em;
		width: 100%;
		max-height: 91vh;
		overflow-x: auto;
		overflow-y: auto;
		z-index: 3;
	}

	.path {
		margin-left: 1vw;
		margin-right: 1vw;
		color: rgb(var(--color-text-primary));
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.info {
		margin-right: 20px;
		display: flex;
		align-items: flex-end;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.info h2 {
		margin-left: 1em;
		color: rgb(var(--color-surface-100));
	}

	table {
		table-layout: auto;
		text-align: left;
		border-collapse: collapse;
		border-spacing: 1px;
		font-size: 13px;
		color: rgb(var(--color-surface-100));
		width: 100%;
		background-color: rgb(var(--color-surface-800));
	}

	td,
	th {
		padding: 10px 20px;
		border-bottom: 1px solid rgb(var(--color-surface-500));
		color: rgb(var(--color-surface-100));
		font-size: 13px;
		white-space: nowrap;
		max-width: 200px;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	input {
		border-color: rgb(var(--color-surface-500));
		outline: none;
	}

	.sticky-header {
		position: sticky;
		top: 0;
		z-index: 1;
		font-size: 13px;
		background-color: rgb(var(--color-surface-600));
		color: rgb(var(--color-surface-100));
	}

	.loader {
		width: 100%;
		height: 4.8px;
		display: inline-block;
		position: relative;
		overflow: hidden;
	}

	.loader::after {
		content: '';
		width: 96px;
		height: 4.8px;
		background: rgb(var(--color-warning-500));
		position: absolute;
		top: 0;
		left: 0;
		box-sizing: border-box;
		animation: hit 0.5s ease-in-out infinite alternate;
	}

	@keyframes hit {
		0% {
			left: 0;
			transform: translateX(-1%);
		}
		100% {
			left: 100%;
			transform: translateX(-99%);
		}
	}
</style>
