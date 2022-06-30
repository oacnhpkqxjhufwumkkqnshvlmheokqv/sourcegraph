import { SeriesSortDirection, SeriesSortMode } from '../../../../../../../../../graphql-operations'

import { DrillDownFiltersFormValues } from './DrillDownInsightFilters'
import { parseSeriesDisplayOptions } from './utils'

const TEST_SERIES_DISPLAY_OPTIONS: DrillDownFiltersFormValues['seriesDisplayOptions'] = {
    limit: '10',
    sortOptions: {
        direction: SeriesSortDirection.ASC,
        mode: SeriesSortMode.DATE_ADDED,
    },
}

describe('BackendInsight', () => {
    describe('parseSeriesDisplayOptions', () => {
        it('returns given object when provided complete values', () => {
            const parsed = parseSeriesDisplayOptions(10, TEST_SERIES_DISPLAY_OPTIONS)
            expect(parsed).toEqual(TEST_SERIES_DISPLAY_OPTIONS)
        })
    })
})
