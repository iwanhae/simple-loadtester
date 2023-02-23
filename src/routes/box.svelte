<script lang="ts">
	import type { ResponseResult } from './types';

	export let result: ResponseResult;
	export let maxDur = 0.0;
	export let minDur = Number.MAX_VALUE;
	export let selected = false;
</script>

<div
	class="box"
	class:failed={result.isFailed}
	class:normal={!result.isFailed && maxDur !== result.duration && result.duration !== undefined}
	class:special={!result.isFailed && maxDur <= (result.duration || 0)}
	style:opacity={maxDur > minDur
		? 1.0 - ((result.duration || 0.0) - minDur) / (maxDur - minDur)
		: 1.0}
	class:selected
/>

<style>
	.box {
		width: 20px;
		height: 20px;
		margin: 1px;
		background-color: gray;
		transition: ease 0.5s;
	}
	.failed {
		background-color: red;
		opacity: 1 !important;
	}
	.normal {
		background-color: aqua;
	}
	.special {
		background-color: blue;
		opacity: 1 !important;
	}
	.selected {
		background-color: black;
		transition: none;
	}
</style>
