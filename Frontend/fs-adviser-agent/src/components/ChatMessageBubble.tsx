'use client';

import * as React from 'react';
import { cn } from '@/lib/utils';

interface ChatMessageBubbleProps {
  role: 'user' | 'assistant';
  className?: string;
  children: React.ReactNode;
}

const ChatMessageBubble = React.forwardRef<HTMLDivElement, ChatMessageBubbleProps>(
  ({ role, className, children }, ref) => {
    return (
      <div
        ref={ref}
        className={cn(
          'whitespace-pre-wrap rounded-lg px-4 py-2 mt-2 text-sm max-w-[85%]',
          role === 'user'
            ? 'bg-gradient-to-r from-blue-700 to-gray-600 text-white ml-auto'
            : 'bg-gray-100 text-gray-800 mr-auto',
          className
        )}
      >
        {children}
      </div>
    );
  }
);

ChatMessageBubble.displayName = 'ChatMessageBubble';

export default ChatMessageBubble;
