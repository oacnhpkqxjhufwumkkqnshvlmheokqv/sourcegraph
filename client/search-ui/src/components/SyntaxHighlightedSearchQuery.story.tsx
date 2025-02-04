import { Meta, Story } from '@storybook/react'

import { BrandedStory } from '@sourcegraph/branded/src/components/BrandedStory'
import { Text } from '@sourcegraph/wildcard'

import { SyntaxHighlightedSearchQuery } from './SyntaxHighlightedSearchQuery'

const config: Meta = {
    title: 'search-ui/SyntaxHighlightedSearchQuery',
    parameters: {
        chromatic: { viewports: [480] },
    },
}

export default config

export const SyntaxHighlightedSearchQueryStory: Story = () => (
    <BrandedStory>
        {() => (
            <Text>
                <SyntaxHighlightedSearchQuery query="test AND spec" />
                <br />
                <SyntaxHighlightedSearchQuery query="test or spec repo:sourcegraph" />
                <br />
                <SyntaxHighlightedSearchQuery query="test -lang:ts" />
            </Text>
        )}
    </BrandedStory>
)

SyntaxHighlightedSearchQueryStory.storyName = 'SyntaxHighlightedSearchQuery'
