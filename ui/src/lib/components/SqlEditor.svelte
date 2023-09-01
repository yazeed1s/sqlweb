<script lang="ts">
	import { onMount, afterUpdate, createEventDispatcher} from 'svelte';
	import loader from '@monaco-editor/loader';
	import type * as Monaco from 'monaco-editor/esm/vs/editor/editor.api';

	export let code: string;
	export let bindCode: string;

	let editor: Monaco.editor.IStandaloneCodeEditor;
	let monaco: typeof Monaco;
	let editorContainer: any;

	const customTheme: Monaco.editor.IStandaloneThemeData = {
		base: 'vs-dark',
		inherit: true,
		rules: [],
		colors: {
			'editor.foreground': '#E9B96E',
			'editor.background': '#121015',
			'editorCursor.foreground': '#BABDB6',
			'editor.lineHighlightBackground': '#1B1820',
			'editorLineNumber.foreground': '#BABDB6',
			'editor.selectionBackground': '#2B282F',
			'editor.inactiveSelectionBackground': '#2B282F'
		}
	};

	const dispatch = createEventDispatcher();

	onMount(async () => {
		const monacoEditor = await import('monaco-editor');
		loader.config({ monaco: monacoEditor.default });

		monaco = await loader.init();
		monaco.editor.defineTheme('custom-theme', customTheme);
		editor = monaco.editor.create(editorContainer, {
			value: code,
			language: 'sql',
			glyphMargin: false,
			theme: 'custom-theme',
			fontSize: 20,
			scrollbar: {
				useShadows: false,
				vertical: 'auto',
				horizontal: 'auto',
				arrowSize: 30
			}
		});

		editor.onDidChangeModelContent(() => {
			const value = editor.getValue();
			bindCode = value;
			dispatch('input', value);
		});
	});

	afterUpdate(() => {
		if (editor && code !== editor.getValue()) {
			editor.setValue(code);
		}
	});

</script>

<div style="width: 100%;">
	<div class="container" bind:this={editorContainer} />
</div>

<style>
	.container {
		position: relative;
		width: 100%;
		height: 60vh;
	}
</style>
