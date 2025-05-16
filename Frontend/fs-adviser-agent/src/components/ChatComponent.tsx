"use client";

import React, { useEffect, useRef } from "react";
import { useChat } from "@ai-sdk/react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { Send } from "lucide-react";
import { Message } from "ai";

type Props = {
  chatId: number;
};

const ChatComponent = ({ chatId }: Props) => {
  const containerRef = useRef<HTMLDivElement | null>(null);

  const { input, handleInputChange, handleSubmit, messages, isLoading } = useChat({
    api: "/api/chat",
    body: { chatId },
  });

  useEffect(() => {
    const el = containerRef.current;
    if (el) {
      el.scrollTo({ top: el.scrollHeight, behavior: "smooth" });
    }
  }, [messages]);

  return (
    <div className="relative max-h-screen flex flex-col h-full">
      {/* Header */}
      <div className="p-4 bg-gradient-to-r from-blue-700 to-gray-600 text-white sticky top-0 z-10">
        <h2 className="text-xl font-semibold">Conversation</h2>
      </div>

      {/* Messages */}
      <div
        ref={containerRef}
        className="flex-1 overflow-y-auto px-4 py-2 space-y-4 bg-white"
      >
        {messages.map((msg: Message) => (
          <div
            key={msg.id}
            className={`max-w-[75%] p-3 rounded-lg shadow-md text-sm ${
              msg.role === "user"
                ? "ml-auto bg-blue-100 text-right"
                : "mr-auto bg-gray-100"
            }`}
          >
            <div className="whitespace-pre-wrap">{msg.content}</div>
          </div>
        ))}

        {isLoading && (
          <div className="flex items-center justify-center min-h-screen bg-gray-50">
            <div className="text-center">
              <div className="w-16 h-16 border-4 border-t-blue-500 border-gray-200 rounded-full animate-spin mx-auto mb-4"></div>
              <p className="text-gray-600 text-lg font-medium">Loading...</p>
            </div>
          </div>
        )}
      </div>

      {/* Input */}
      <form
        onSubmit={handleSubmit}
        className="p-4 border-t border-gray-200 bg-white sticky bottom-0 z-10"
      >
        <div className="flex gap-2">
          <Input
            placeholder="Ask anything about the market..."
            value={input}
            onChange={handleInputChange}
            className="flex-1"
          />
          <Button type="submit" className="bg-gradient-to-r from-blue-700 to-gray-600">
            <Send className="h-4 w-4" />
          </Button>
        </div>
      </form>
    </div>
  );
};

export default ChatComponent;
