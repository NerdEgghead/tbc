import { ActionMetrics, SimResult, SimResultFilter } from '/tbc/core/proto_utils/sim_result.js';
import { sum } from '/tbc/core/utils.js';

import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

declare var $: any;
declare var tippy: any;

export class MeleeMetrics extends ResultComponent {
	private readonly tableElem: HTMLTableSectionElement;
	private readonly bodyElem: HTMLTableSectionElement;

  constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'cast-metrics-root';
    super(config);

		this.rootElem.innerHTML = `
		<table class="metrics-table tablesorter">
			<thead class="metrics-table-header">
				<tr class="metrics-table-header-row">
					<th class="metrics-table-header-cell"><span>Name</span></th>
					<th class="metrics-table-header-cell"><span>DPS</span></th>
					<th class="metrics-table-header-cell"><span>Casts</span></th>
					<th class="metrics-table-header-cell"><span>Avg Cast</span></th>
					<th class="metrics-table-header-cell"><span>Hits</span></th>
					<th class="metrics-table-header-cell"><span>Avg Hit</span></th>
					<th class="metrics-table-header-cell"><span>Crit %</span></th>
					<th class="metrics-table-header-cell"><span>Miss %</span></th>
					<th class="metrics-table-header-cell"><span>Dodge %</span></th>
					<th class="metrics-table-header-cell"><span>Glance %</span></th>
				</tr>
			</thead>
			<tbody class="metrics-table-body">
			</tbody>
		</table>
		`;

		this.tableElem = this.rootElem.getElementsByClassName('metrics-table')[0] as HTMLTableSectionElement;
		this.bodyElem = this.rootElem.getElementsByClassName('metrics-table-body')[0] as HTMLTableSectionElement;

		const headerElems = Array.from(this.tableElem.querySelectorAll('th'));

		// DPS
		tippy(headerElems[1], {
			'content': 'Damage / Encounter Duration',
			'allowHTML': true,
		});

		// Casts
		tippy(headerElems[2], {
			'content': 'Casts',
			'allowHTML': true,
		});

		// Avg Cast
		tippy(headerElems[3], {
			'content': 'Damage / Casts',
			'allowHTML': true,
		});

		// Hits
		tippy(headerElems[4], {
			'content': 'Hits + Crits + Glances + Blocks',
			'allowHTML': true,
		});

		// Avg Hit
		tippy(headerElems[5], {
			'content': 'Damage / (Hits + Crits + Glances + Blocks)',
			'allowHTML': true,
		});

		// Crit %
		tippy(headerElems[6], {
			'content': 'Crits / Swings',
			'allowHTML': true,
		});

		// Miss %
		tippy(headerElems[7], {
			'content': 'Misses / Swings',
			'allowHTML': true,
		});

		// Dodge %
		tippy(headerElems[8], {
			'content': 'Dodges / Swings',
			'allowHTML': true,
		});

		// Glance %
		tippy(headerElems[9], {
			'content': 'Glances / Swings',
			'allowHTML': true,
		});

		$(this.tableElem).tablesorter({ sortList: [[1, 1]] });
	}

	onSimResult(resultData: SimResultData) {
		this.bodyElem.textContent = '';

		const players = resultData.result.getPlayers(resultData.filter);
		if (players.length != 1) {
			return;
		}
		const player = players[0];

		const addRow = (meleeMetric: ActionMetrics, isChildMetric: boolean, rowNamePrefix: string, rowName: string): HTMLElement => {
			const fullName = rowNamePrefix ? rowNamePrefix + ' - ' + rowName : rowName;

			const rowElem = document.createElement('tr');
			if (isChildMetric) {
				rowElem.classList.add('child-metric');
			}
			this.bodyElem.appendChild(rowElem);

			const nameCellElem = document.createElement('td');
			rowElem.appendChild(nameCellElem);
			nameCellElem.innerHTML = `
			<a class="metrics-action-icon"></a>
			<span class="metrics-action-name">${fullName}</span>
			<span class="expand-toggle fa fa-caret-down"></span>
			<span class="expand-toggle fa fa-caret-right"></span>
			`;

			const iconElem = nameCellElem.getElementsByClassName('metrics-action-icon')[0] as HTMLAnchorElement;
			meleeMetric.actionId.setBackgroundAndHref(iconElem);

			const addCell = (value: string | number): HTMLElement => {
				const cellElem = document.createElement('td');
				cellElem.textContent = String(value);
				rowElem.appendChild(cellElem);
				return cellElem;
			};

			addCell(meleeMetric.dps.toFixed(1)); // DPS
			addCell(meleeMetric.casts.toFixed(1)); // Casts
			addCell(meleeMetric.avgCast.toFixed(1)); // Avg Cast
			addCell(meleeMetric.landedHits.toFixed(1)); // Hits
			addCell(meleeMetric.avgHit.toFixed(1)); // Avg Hit
			addCell(meleeMetric.critPercent.toFixed(2) + ' %'); // Crit %
			addCell(meleeMetric.missPercent.toFixed(2) + ' %'); // Miss %
			addCell(meleeMetric.dodgePercent.toFixed(2) + ' %'); // Dodge %
			addCell(meleeMetric.glancePercent.toFixed(2) + ' %'); // Glance %
			return rowElem;
		};
		const addGroup = (meleeGroup: Array<ActionMetrics>, namePrefix: string) => {
			if (meleeGroup.length == 1) {
				addRow(meleeGroup[0], false, namePrefix, meleeGroup[0].name);
				return;
			}

			const mergedMetrics = ActionMetrics.merge(meleeGroup, true);
			const parentRow = addRow(mergedMetrics, false, '', namePrefix || mergedMetrics.name);
			const childRows = meleeGroup.map(meleeMetric => addRow(meleeMetric, true, namePrefix, meleeMetric.name));
			const defaultDisplay = childRows[0].style.display;

			let expand = true;
			parentRow.classList.add('parent-metric', 'expand');
			parentRow.addEventListener('click', event => {
				expand = !expand;
				const newDisplayValue = expand ? defaultDisplay : 'none';
				childRows.forEach(row => row.style.display = newDisplayValue);
				if (expand) {
					parentRow.classList.add('expand');
				} else {
					parentRow.classList.remove('expand');
				}
			});
		};

		const meleeMetrics = player.getMeleeActions();
		const meleeGroups = ActionMetrics.groupById(meleeMetrics);

		if (meleeMetrics.length == 0) {
			this.rootElem.classList.add('empty');
			return
		} else {
			this.rootElem.classList.remove('empty');
		}

		meleeGroups.forEach(meleeGroup => addGroup(meleeGroup, ''));
		player.pets.forEach(pet => addGroup(pet.getMeleeActions(), pet.name));

		$(this.tableElem).trigger('update');
	}
}
