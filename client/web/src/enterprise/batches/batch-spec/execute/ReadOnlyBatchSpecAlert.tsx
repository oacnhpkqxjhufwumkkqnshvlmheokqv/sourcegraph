import React from 'react'

import { Alert, AlertProps, H4 } from '@sourcegraph/wildcard'

interface ReadOnlyBatchSpecAlertProps {
    className?: string
    variant: AlertProps['variant']
    header: string
    message: React.ReactNode
}

export const ReadOnlyBatchSpecAlert: React.FunctionComponent<React.PropsWithChildren<ReadOnlyBatchSpecAlertProps>> = ({
    className,
    children,
    variant,
    header,
    message,
}) => (
    <Alert variant={variant} className={className}>
        <div className="flex-grow-1 pr-3">
            <H4>{header}</H4>
            {message}
        </div>
        {children}
    </Alert>
)
