<script lang="ts">
	import SqlEditor from './SqlEditor.svelte';
	import { connectedStore } from '$lib/store';
	import { slide } from 'svelte/transition';
	import { endpoints, httpClient } from '../../services/api';

	let connected: boolean = false;
	let req = { query: 'SELECT * FROM test_table;' };
	let isError: boolean = false;
	let isSuccess: boolean = false;
	let isTable: boolean = false;
	let error: string = '';
	let msg: string = '';
	let header: string[] = [];
	let rows: Array<any> = [];

	connectedStore.subscribe((value) => {
		connected = value;
	});

	const sendRequest = async (): Promise<void> => {
		const res = await httpClient(endpoints.execute, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(req)
		});
		const json = await res.json();
		console.log(json);
		res.ok ? handleSuccessResponse(json) : handleErrorResponse(json);
	};

	const handleInput = (event: { detail: string }): void => {
		req.query = event.detail;
	};

	const runQuery = async (): Promise<void> => {
		await sendRequest();
	};

	const handleSuccessResponse = (json: any): void => {
		// TODO: fix this shit (detect if the resp is a table or not)
		if (json.data.result.data) {
			isTable = true;
			header = Object.keys(json.data.result.data[0]);
			rows = json.data.result.data.map((row: any) => Object.values(row));
		}
		isSuccess = true;
		isError = false;
		msg = '';
		msg = json.data.result.message;
	};

	const handleErrorResponse = (json: any): void => {
		// TODO: handle this unnecessary values resetting
		isSuccess = false;
		isError = true;
		error = '';
		error = json.error;
		msg = '';
		msg = json.message;
	};
</script>

{#if connected}
	<body>
		<div class="content">
			<main>
				<SqlEditor bind:code={req.query} bind:bindCode={req.query} on:input={handleInput} />
				{#if isTable}
					<div class="table-container">
						<table>
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
								{#each rows as row}
									<tr>
										{#each row as value}
											<td>{value}</td>
										{/each}
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
				{#if isError}
					<div class="error" role="alert">
						<p class="font-error">Error: {error}</p>
						<p class="msg">{msg}</p>
					</div>
				{/if}
			</main>
		</div>
		{#if isSuccess}
			<div class="success" role="alert">
				<p class="font">Success: {msg}</p>
			</div>
		{/if}
		<button class="run-button" on:click={() => runQuery()}>Run Query</button>
	</body>
{/if}

<style>
	.table-container {
		margin-top: 1em;
		width: 100%;
		padding-right: 1vw;
		padding-left: 1vw;
		overflow-x: auto;
		overflow-y: auto;
		max-height: 73vh;
	}
	table {
		table-layout: auto;
		text-align: left;
		border-collapse: collapse;
		border-spacing: 1;
		font-size: 13px;
		color: #313c45;
		width: 100%;
		max-height: 80.5vh;
		background-color: rgb(var(--color-surface-800));
	}

	td,
	th {
		padding: 4px 20px;
		border-bottom: 1px solid rgb(var(--color-surface-500));
		color: rgb(var(--color-surface-100));
		font-size: 13px;
	}

	.sticky-header {
		position: sticky;
		top: 0;
		z-index: 1;
		font-size: 13px;
		background-color: rgb(var(--color-surface-600));
		color: rgb(var(--color-surface-100));
	}

	th {
		color: rgb(var(--color-surface-100));
		font-size: 13px;
	}

	.error {
		background-color: rgb(var(--color-surface-500));
		border-left-width: 4px;
		border-color: rgb(var(--color-error-500));
		padding: 0.7rem;
	}

	.error .font-error {
		font-weight: 700;
		color: rgb(var(--color-error-600));
	}

	.error .msg {
		font-weight: 500;
		color: rgb(var(--color-surface-800));
	}

	.success {
		background-color: rgb(var(--color-surface-500));
		border-left-width: 4px;
		border-color: rgb(var(--color-success-500));
		padding: 0.7rem;
	}
	.success .font {
		font-weight: 700;
		color: rgb(var(--color-surface-800));
	}

	.run-button {
		margin-bottom: 2vh;
		margin-left: 1vw;
		background-color: rgb(var(--color-success-500));
		color: rgb(var(--color-surface-800));
		border-radius: 10px;
		width: 80px;
		height: 30px;
	}

	.run-button:hover {
		background-color: rgb(var(--color-success-500));
	}

	* {
		margin: 0px;
		box-sizing: border-box;
	}

	body {
		height: 100%;
		width: 100%;
		text-align: start;
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
		padding: 0em 0 0em 0;
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
