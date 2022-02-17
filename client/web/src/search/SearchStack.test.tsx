import { cleanup, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import React from 'react'

import { SearchPatternType } from '@sourcegraph/shared/src/graphql-operations'
import { renderWithBrandedContext, RenderWithBrandedContextResult } from '@sourcegraph/shared/src/testing'

import { useExperimentalFeatures, useSearchStackState } from '../stores'
import { SearchStackEntry } from '../stores/searchStack'

import { SearchStack } from './SearchStack'

describe('Search Stack', () => {
    const renderSearchStack = (props?: Partial<{ initialOpen: boolean }>): RenderWithBrandedContextResult =>
        renderWithBrandedContext(<SearchStack {...props} />)

    afterEach(cleanup)

    const mockEntries: SearchStackEntry[] = [
        { id: 0, type: 'search', query: 'TODO', caseSensitive: false, patternType: SearchPatternType.literal },
        { id: 1, type: 'file', path: 'path/to/file', repo: 'test', revision: 'master', lineRange: null },
    ]

    describe('inital state', () => {
        it('does not render anything if feature is disabled', () => {
            useExperimentalFeatures.setState({ enableSearchStack: false })
            useSearchStackState.setState({ addableEntry: mockEntries[0] })

            renderSearchStack()

            expect(screen.queryByRole('button', { name: 'Add search' })).not.toBeInTheDocument()
        })

        it('shows the add button if an entry can be added', () => {
            useExperimentalFeatures.setState({ enableSearchStack: true })
            useSearchStackState.setState({ canRestoreSession: true, addableEntry: mockEntries[0] })

            expect(renderSearchStack().asFragment()).toMatchSnapshot()
        })

        it('shows the top of the stack if entries exist', () => {
            useExperimentalFeatures.setState({ enableSearchStack: true })
            useSearchStackState.setState({ canRestoreSession: true, entries: mockEntries })

            expect(renderSearchStack().asFragment()).toMatchSnapshot()
        })
    })

    describe('restore previous session', () => {
        beforeEach(() => {
            useExperimentalFeatures.setState({ enableSearchStack: true })
        })

        it('restores the previous session', () => {
            useSearchStackState.setState({
                entries: [],
                previousEntries: mockEntries,
                canRestoreSession: true,
                addableEntry: mockEntries[0],
            })
            renderSearchStack()
            userEvent.click(screen.getByRole('button', { name: 'Open search session' }))

            userEvent.click(screen.getByRole('button', { name: 'Restore previous session' }))
            expect(useSearchStackState.getState().entries).toEqual(mockEntries)
        })
    })

    describe('with entries', () => {
        beforeEach(() => {
            useExperimentalFeatures.setState({ enableSearchStack: true })
            useSearchStackState.setState({
                entries: [
                    {
                        id: 0,
                        type: 'search',
                        query: 'TODO',
                        caseSensitive: false,
                        patternType: SearchPatternType.literal,
                    },
                    { id: 1, type: 'file', path: 'path/to/file', repo: 'test', revision: 'master', lineRange: null },
                ],
            })
        })

        it('opens and closes', () => {
            renderSearchStack()

            userEvent.click(screen.getByRole('button', { name: 'Open search session' }))

            const closeButtons = screen.queryAllByRole('button', { name: 'Close search session' })
            expect(closeButtons).toHaveLength(2)

            userEvent.click(closeButtons[0])
            expect(screen.queryByRole('button', { name: 'Open search session' })).toBeInTheDocument()
        })

        it('redirects to entries', () => {
            renderSearchStack()
            userEvent.click(screen.getByRole('button', { name: 'Open search session' }))

            const entryLinks = screen.queryAllByRole('link')

            // Entries are in reverse order
            expect(entryLinks[0]).toHaveAttribute('href', '/test@master/-/blob/path/to/file')
            expect(entryLinks[1]).toHaveAttribute('href', '/search?q=TODO&patternType=literal')
        })

        it('creates notebooks', () => {
            const result = renderSearchStack()

            userEvent.click(screen.getByRole('button', { name: 'Open search session' }))
            userEvent.click(screen.getByRole('button', { name: 'Create Notebook' }))

            expect(result.history.location.pathname).toMatchInlineSnapshot('"/notebooks/new"')
            expect(result.history.location.hash).toMatchInlineSnapshot(
                '"#query:TODO,file:http%3A%2F%2Flocalhost%2Ftest%40master%2F-%2Fblob%2Fpath%2Fto%2Ffile"'
            )
        })

        it('allows to delete entries', () => {
            renderSearchStack()
            userEvent.click(screen.getByRole('button', { name: 'Open search session' }))

            userEvent.click(screen.getAllByRole('button', { name: 'Remove entry' })[0])
            const entryLinks = screen.queryByRole('link')
            expect(entryLinks).toBeInTheDocument()
        })

        it('opens the text annotation aria', () => {
            renderSearchStack()
            userEvent.click(screen.getByRole('button', { name: 'Open search session' }))

            userEvent.click(screen.getAllByRole('button', { name: 'Add annotation' })[0])
            expect(screen.queryByPlaceholderText('Type to add annotation...')).toBeInTheDocument()
        })
    })
})
