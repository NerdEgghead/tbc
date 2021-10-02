import { TypedEvent } from '../core/typed_event.js';
import { DpsHistogram } from './dps_histogram.js';
import { DpsResult } from './dps_result.js';
import { PercentOom } from './percent_oom.js';
const urlParams = new URLSearchParams(window.location.search);
if (urlParams.has('mainBgColor')) {
    document.body.style.setProperty('--main-bg-color', urlParams.get('mainBgColor'));
}
if (urlParams.has('mainTextColor')) {
    document.body.style.setProperty('--main-text-color', urlParams.get('mainTextColor'));
}
const colorSettings = {
    mainTextColor: document.body.style.getPropertyValue('--main-text-color'),
};
Chart.defaults.color = colorSettings.mainTextColor;
const layoutHTML = `
<div class="dr-root">
	<div class="dr-row topline-results">
	</div>
	<div class="dr-row dps-histogram">
	</div>
</div>
`;
const resultsEmitter = new TypedEvent();
window.addEventListener('message', event => {
    // Null indicates pending results
    const data = event.data;
    resultsEmitter.emit(event.data);
});
document.body.innerHTML = layoutHTML;
const toplineResultsDiv = document.body.getElementsByClassName('topline-results')[0];
const dpsResult = new DpsResult({ parent: toplineResultsDiv, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const percentOom = new PercentOom({ parent: toplineResultsDiv, resultsEmitter: resultsEmitter, colorSettings: colorSettings });
const dpsHistogram = new DpsHistogram({ parent: document.body.getElementsByClassName('dps-histogram')[0], resultsEmitter: resultsEmitter, colorSettings: colorSettings });
