import { FunctionComponent } from 'react'

import BarChartIcon from 'mdi-react/BarChartIcon'

import { TelemetryProps } from '@sourcegraph/shared/src/telemetry/telemetryService'
import { Link, PageHeader, Text, useLocalStorage } from '@sourcegraph/wildcard'

import { PageTitle } from '../../../../../../components/PageTitle'
import { CodeInsightsPage } from '../../../../components/code-insights-page/CodeInsightsPage'
import { FormChangeEvent } from '../../../../components/form/hooks/useForm'

import { ComputeInsightCreationContent } from './components/ComputeInsightCreationContent'
import { CreateComputeInsightFormFields } from './components/types'

export interface InsightCreateEvent {
    // TODO: It will be improved in https://github.com/sourcegraph/sourcegraph/issues/37965
    insight: any
}

interface ComputeInsightCreationPageProps extends TelemetryProps {
    onInsightCreateRequest: (event: InsightCreateEvent) => Promise<unknown>
    onSuccessfulCreation: () => void
    onCancel: () => void
}

export const ComputeInsightCreationPage: FunctionComponent<ComputeInsightCreationPageProps> = props => {
    const [initialFormValues, setInitialFormValues] = useLocalStorage<CreateComputeInsightFormFields | undefined>(
        'insights.compute-creation-ui',
        undefined
    )

    const handleChange = (event: FormChangeEvent<CreateComputeInsightFormFields>): void => {
        setInitialFormValues(event.values)
    }

    const handleSubmit = (): void => {
        // TODO: It will be implemented in https://github.com/sourcegraph/sourcegraph/issues/37965
    }

    const handleCancel = (): void => {
        // TODO: It will be implemented in https://github.com/sourcegraph/sourcegraph/issues/37965
    }

    return (
        <CodeInsightsPage className="col-10">
            <PageTitle title="Create compute insight - Code Insights" />

            <PageHeader
                className="mb-5"
                path={[{ icon: BarChartIcon }, { text: 'Create code insight' }]}
                description={
                    <Text>
                        Type: <b>Map results</b> | <Link to="/insights/create">Change type</Link>
                    </Text>
                }
            />

            <ComputeInsightCreationContent
                initialValue={initialFormValues}
                data-testid="search-insight-create-page-content"
                className="pb-5"
                onChange={handleChange}
                onSubmit={handleSubmit}
                onCancel={handleCancel}
            />
        </CodeInsightsPage>
    )
}
