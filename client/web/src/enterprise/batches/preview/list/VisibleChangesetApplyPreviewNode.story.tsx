import { Meta, DecoratorFn, Story } from '@storybook/react'
import classNames from 'classnames'
import { of } from 'rxjs'

import { WebStory } from '../../../../components/WebStory'
import {
    ChangesetSpecOperation,
    ChangesetState,
    ChangesetSpecType,
    VisibleChangesetApplyPreviewFields,
} from '../../../../graphql-operations'

import { baseChangesetSpec, testRepo } from './testData'
import { VisibleChangesetApplyPreviewNode } from './VisibleChangesetApplyPreviewNode'

import styles from './PreviewList.module.scss'

const decorator: DecoratorFn = story => (
    <div className={classNames(styles.previewListGrid, 'p-3 container')}>{story()}</div>
)

const config: Meta = {
    title: 'web/batches/preview/VisibleChangesetApplyPreviewNode',
    decorators: [decorator],
}

export default config

const visibleChangesetApplyPreviewNodeStories = (
    publicationStateSet: boolean
): Record<string, VisibleChangesetApplyPreviewFields> => ({
    ImportChangeset: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.IMPORT],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsAttach',
            changesetSpec: {
                __typename: 'VisibleChangesetSpec',
                id: 'someidv1',
                type: ChangesetSpecType.EXISTING,
                description: {
                    __typename: 'ExistingChangesetReference',
                    baseRepository: { name: 'github.com/sourcegraph/testrepo', url: 'https://test.test/repo' },
                    externalID: '123',
                },
                forkTarget: null,
            },
        },
    },
    CreateChangesetPublished: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.PUSH, ChangesetSpecOperation.PUBLISH],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsAttach',
            changesetSpec: baseChangesetSpec(1, publicationStateSet ? true : null),
        },
    },
    CreateChangesetDraft: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.PUSH, ChangesetSpecOperation.PUBLISH_DRAFT],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsAttach',
            changesetSpec: baseChangesetSpec(2, publicationStateSet ? 'draft' : null),
        },
    },
    CreateChangesetNotPublished: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsAttach',
            changesetSpec: baseChangesetSpec(3, publicationStateSet ? false : null),
        },
    },
    UpdateChangesetTitle: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.UPDATE],
        delta: {
            titleChanged: true,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(4, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'the old title',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: null,
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    UpdateChangesetBody: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.UPDATE],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: true,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(5, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'the old title',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: null,
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    UndraftChangeset: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.UNDRAFT],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(6, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'Le draft changeset',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: null,
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    ReopenChangeset: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.REOPEN, ChangesetSpecOperation.UPDATE],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(7, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'Le closed changeset',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: null,
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    CloseChangeset: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.CLOSE, ChangesetSpecOperation.DETACH],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsDetach',
            changeset: {
                id: '123123',
                title: 'Le open changeset',
                state: ChangesetState.OPEN,
                repository: testRepo,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                diffStat: {
                    added: 2,
                    changed: 8,
                    deleted: 10,
                },
            },
        },
    },
    DetachChangeset: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.DETACH],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsDetach',
            changeset: {
                id: '123123',
                title: 'Le open changeset',
                state: ChangesetState.OPEN,
                repository: testRepo,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                diffStat: {
                    added: 2,
                    changed: 8,
                    deleted: 10,
                },
            },
        },
    },
    ChangeBaseRef: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.UPDATE],
        delta: {
            titleChanged: false,
            baseRefChanged: true,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(8, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'Change base ref',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: null,
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    ChangeDiff: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.UPDATE],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: true,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(9, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'Change base ref',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'master',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: null,
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    UpdateCommitMessage: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.PUSH],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: true,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(10, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'the old title',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: 'Current commit message',
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Alice',
                    email: 'alice@email.test',
                    user: {
                        displayName: 'Alice',
                        url: '/users/alice',
                        username: 'alice',
                    },
                },
            },
        },
    },
    UpdateCommitAuthor: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.PUSH],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: true,
            authorNameChanged: true,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsUpdate',
            changesetSpec: baseChangesetSpec(11, publicationStateSet ? true : null),
            changeset: {
                id: '123123',
                title: 'the old title',
                state: ChangesetState.OPEN,
                externalID: '123',
                externalURL: {
                    url: 'http://test.test/123',
                },
                currentSpec: {
                    description: {
                        __typename: 'GitBranchChangesetDescription',
                        baseRef: 'main',
                        body: 'body',
                        commits: [
                            {
                                subject: 'Abc',
                                body: 'Current commit message',
                                author: {
                                    avatarURL: null,
                                    displayName: 'alice',
                                    email: 'alice@sourcegraph.test',
                                    user: null,
                                },
                            },
                        ],
                        title: 'Title',
                    },
                },
                author: {
                    displayName: 'Bob',
                    email: 'bob@email.test',
                    user: {
                        displayName: 'Bob',
                        url: '/users/bob',
                        username: 'bob',
                    },
                },
            },
        },
    },
    ForkedRepo: {
        __typename: 'VisibleChangesetApplyPreview',
        operations: [ChangesetSpecOperation.PUSH, ChangesetSpecOperation.PUBLISH],
        delta: {
            titleChanged: false,
            baseRefChanged: false,
            diffChanged: false,
            bodyChanged: false,
            authorEmailChanged: false,
            authorNameChanged: false,
            commitMessageChanged: false,
        },
        targets: {
            __typename: 'VisibleApplyPreviewTargetsAttach',
            changesetSpec: baseChangesetSpec(12, publicationStateSet ? true : null, {
                forkTarget: { pushUser: true, namespace: null },
                description: {
                    __typename: 'GitBranchChangesetDescription',
                    baseRepository: testRepo,
                    baseRef: 'master',
                    headRef: 'cool-branch',
                    body: 'Body text',
                    commits: [
                        {
                            subject: 'This is the first line of the commit message',
                            body: `And the more explanatory body. And the more explanatory body.
And the more explanatory body. And the more explanatory body.
And the more explanatory body. And the more explanatory body.
And the more explanatory body. And the more explanatory body. And the more explanatory body.
And the more explanatory body. And the more explanatory body. And the more explanatory body.`,
                            author: {
                                avatarURL: null,
                                displayName: 'john',
                                email: 'john@test.not',
                                user: { displayName: 'lejohn', url: '/users/lejohn', username: 'john' },
                            },
                        },
                    ],
                    diffStat: {
                        __typename: 'DiffStat',
                        added: 10,
                        changed: 8,
                        deleted: 2,
                    },
                    title: 'Add prettier to forked repository',
                    published: publicationStateSet,
                },
            }),
        },
    },
})

const queryEmptyFileDiffs = () => of({ totalCount: 0, pageInfo: { endCursor: null, hasNextPage: false }, nodes: [] })

const stories = visibleChangesetApplyPreviewNodeStories(true)

export const ImportChangeset: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.ImportChangeset}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

ImportChangeset.storyName = 'Import changeset'

export const CreateChangesetPublished: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.CreateChangesetPublished}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

CreateChangesetPublished.storyName = 'Create changeset published'

export const CreateChangesetDraft: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.CreateChangesetDraft}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

CreateChangesetDraft.storyName = 'Create changeset draft'

export const CreateChangesetNotPublished: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.CreateChangesetNotPublished}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

CreateChangesetNotPublished.storyName = 'Create changeset not published'

export const UpdateChangesetTitle: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.UpdateChangesetTitle}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

UpdateChangesetTitle.storyName = 'Update changeset tittle'

export const UpdateChangesetBody: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.UpdateChangesetBody}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

UpdateChangesetBody.storyName = 'Update changeset body'

export const UndraftChangeset: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.UndraftChangeset}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

UndraftChangeset.storyName = 'Undraft changeset'

export const ReopenChangeset: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.ReopenChangeset}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

ReopenChangeset.storyName = 'Reopen changeset'

export const CloseChangeset: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.CloseChangeset}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

CloseChangeset.storyName = 'Close changeset'

export const DetachChangeset: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.DetachChangeset}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

DetachChangeset.storyName = 'Detach changeset'

export const ChangeBaseRef: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.ChangeBaseRef}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

ChangeBaseRef.storyName = 'Change base ref'

export const ChangeDiff: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.ChangeDiff}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

ChangeDiff.storyName = 'Change diff'

export const UpdateCommitMessage: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.UpdateCommitMessage}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

UpdateCommitMessage.storyName = 'Update commit message'

export const UpdateCommitAuthor: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.UpdateCommitAuthor}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

UpdateCommitAuthor.storyName = 'Update commit author'

export const ForkedRepo: Story = () => (
    <WebStory>
        {props => (
            <VisibleChangesetApplyPreviewNode
                {...props}
                node={stories.ForkedRepo}
                authenticatedUser={{
                    url: '/users/alice',
                    displayName: 'Alice',
                    username: 'alice',
                    email: 'alice@email.test',
                }}
                queryChangesetSpecFileDiffs={queryEmptyFileDiffs}
            />
        )}
    </WebStory>
)

ForkedRepo.storyName = 'Forked repo'
