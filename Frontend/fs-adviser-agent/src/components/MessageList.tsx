"use client";
import { Message } from "ai";
import React from "react";
import { cn } from "@/lib/utils";

type MessageListProps = {
  messages: Message[];
  isLoading?: boolean;
};

const MessageList: React.FC<MessageListProps> = ({ messages, isLoading }) => {
  if (isLoading) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        Loading messages...
      </div>
    );
  }

  if (!messages.length) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        No messages yet.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4 px-4 py-2">
      {messages.map((msg) => (
        <div
          key={msg.id}
          className={cn(
            "p-3 rounded-md max-w-xl shadow text-sm whitespace-pre-wrap",
            msg.role === "user"
              ? "bg-gradient-to-r from-blue-600 to-gray-600 text-white self-end"
              : "bg-gray-100 text-gray-800 self-start"
          )}
        >
          {msg.content}
        </div>
      ))}
    </div>
  );
};

export default MessageList;
