<script lang="ts">
	import type { ResponseResult } from './types';

	export let result: ResponseResult;
	export let maxDur = 0.0;
	export let minDur = Number.MAX_VALUE;
	export let selected = false;

	export let useOpacity = false;
	export let useTransition = false;
</script>

<div
	class="outter"
	class:httpFailed={!result.isFailed && 299 < (result.status || 0)}
	class:special={!result.isFailed && maxDur <= (result.duration || 0)}
>
	<div
		class="box"
		class:networkFailed={result.isFailed}
		style:opacity={useOpacity
			? maxDur > minDur
				? 1.0 - ((result.duration || 0.0) - minDur) / (maxDur - minDur)
				: 1.0
			: 1.0}
		style:background-color={result.color || 'gray'}
		class:selected
		class:transition={useTransition}
	/>
</div>

<style>
	.outter {
		width: 18px;
		height: 18px;
		border: solid white 2px;
	}
	.box {
		width: 18px;
		height: 18px;
	}
	.transition {
		transition: ease 0.5s;
	}
	.networkFailed {
		background-color: red !important;
		opacity: 1 !important;
	}
	.httpFailed {
		border-color: red;
	}
	.special {
		border-color: blue;
	}
	.selected {
		background-color: black !important;
		transition: none;
	}
</style>
