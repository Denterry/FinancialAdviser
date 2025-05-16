'use client';

import { useEffect, useState } from 'react';
import ChatComponent from '@/components/ChatComponent';
import ChatSidebar from '@/components/ChatSideBar';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';

export default function ChatPage({ params }: { params: { chatId: string } }) {
  const [selectedChatId, setSelectedChatId] = useState<number | null>(null);

  const { data: chats = [], isLoading } = useQuery({
    queryKey: ['chats'],
    queryFn: async () => {
      const res = await axios.get('/api/get-chats');
      return res.data;
    },
  });

  useEffect(() => {
    if (chats.length > 0 && selectedChatId === null) {
      setSelectedChatId(parseInt(params.chatId));
    }
  }, [chats, params.chatId, selectedChatId]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-t-blue-500 border-gray-200 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600 text-lg font-medium">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex max-h-screen overflow-hidden">
      {/* Sidebar */}
      <div className="flex-[1] max-w-xs border-r border-gray-200">
        <ChatSidebar
          chats={chats}
          selectedChatId={selectedChatId}
          onSelect={setSelectedChatId}
        />
      </div>

      {/* Main chat component */}
      <div className="flex-[3] overflow-y-auto h-screen">
        {selectedChatId !== null && <ChatComponent chatId={selectedChatId} />}
      </div>
    </div>
  );
}
