'use client';

import React from 'react';
import Link from 'next/link';
import { Button } from './ui/button';
import { Home, MessageCircle, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';

type Chat = {
  id: number;
  title: string;
};

type Props = {
  chats: Chat[];
  selectedChatId: number | null;
  onSelect: (chatId: number) => void;
};

const ChatSideBar: React.FC<Props> = ({ chats, selectedChatId, onSelect }) => {
  return (
    <aside className="w-full h-screen overflow-y-auto bg-gradient-to-b from-blue-900 to-gray-800 text-white p-4">
      <div className="flex flex-col gap-4">
        {/* New Chat Button */}
        <Link href="/chat">
          <Button className="w-full bg-gradient-to-r from-blue-700 to-gray-600 border border-white border-dashed text-white hover:opacity-90">
            <Plus className="mr-2 w-4 h-4" />
            New Chat
          </Button>
        </Link>

        {/* List of Chats */}
        <div className="flex flex-col gap-2 mt-2">
          {chats.length === 0 && (
            <p className="text-sm text-gray-300 italic">No chats yet.</p>
          )}

          {chats.map((chat) => (
            <button
              key={chat.id}
              onClick={() => onSelect(chat.id)}
              className={cn(
                'flex items-center gap-2 p-3 rounded-lg text-sm text-slate-300 transition-colors',
                chat.id === selectedChatId
                  ? 'bg-gradient-to-r from-blue-700 to-gray-600 text-white'
                  : 'hover:text-white hover:bg-gray-700'
              )}
            >
              <MessageCircle className="w-4 h-4" />
              <span className="truncate">{chat.title}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Footer */}
      <div className="absolute bottom-6 left-4">
        <Link href="/">
          <Home className="w-6 h-6 text-gray-300 hover:text-white transition-colors" />
        </Link>
      </div>
    </aside>
  );
};

export default ChatSideBar;
