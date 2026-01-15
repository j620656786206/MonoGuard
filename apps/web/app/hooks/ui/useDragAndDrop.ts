'use client';

import { useState, useCallback, DragEvent } from 'react';

export interface DragAndDropState {
  isDragOver: boolean;
  isDragActive: boolean;
}

export interface DragAndDropActions {
  onDragEnter: (e: DragEvent<HTMLElement>) => void;
  onDragOver: (e: DragEvent<HTMLElement>) => void;
  onDragLeave: (e: DragEvent<HTMLElement>) => void;
  onDrop: (e: DragEvent<HTMLElement>) => void;
  reset: () => void;
}

export interface UseDragAndDropProps {
  onFileDrop: (files: File[]) => void;
  accept?: string[];
}

export interface UseDragAndDropReturn
  extends DragAndDropState,
    DragAndDropActions {}

export const useDragAndDrop = ({
  onFileDrop,
  accept = [],
}: UseDragAndDropProps): UseDragAndDropReturn => {
  const [state, setState] = useState<DragAndDropState>({
    isDragOver: false,
    isDragActive: false,
  });

  const validateFiles = useCallback(
    (files: File[]): File[] => {
      if (accept.length === 0) return files;

      return files.filter((file) => {
        const fileName = file.name.toLowerCase();
        const fileType = file.type.toLowerCase();

        return accept.some((acceptType) => {
          // Check file extension
          if (acceptType.startsWith('.')) {
            return (
              fileName.endsWith(acceptType) ||
              fileName === acceptType.substring(1)
            );
          }

          // Check MIME type
          if (acceptType.includes('/')) {
            return (
              fileType === acceptType ||
              fileType.startsWith(acceptType.replace('*', ''))
            );
          }

          return false;
        });
      });
    },
    [accept]
  );

  const getFilesFromEvent = useCallback(
    (e: DragEvent<HTMLElement>): File[] => {
      const files: File[] = [];

      if (e.dataTransfer?.items) {
        // Use DataTransferItemList interface
        for (let i = 0; i < e.dataTransfer.items.length; i++) {
          const item = e.dataTransfer.items[i];
          if (item.kind === 'file') {
            const file = item.getAsFile();
            if (file) files.push(file);
          }
        }
      } else if (e.dataTransfer?.files) {
        // Use FileList interface
        for (let i = 0; i < e.dataTransfer.files.length; i++) {
          files.push(e.dataTransfer.files[i]);
        }
      }

      return validateFiles(files);
    },
    [validateFiles]
  );

  const onDragEnter = useCallback((e: DragEvent<HTMLElement>) => {
    e.preventDefault();
    e.stopPropagation();

    setState((prev) => ({
      ...prev,
      isDragActive: true,
    }));
  }, []);

  const onDragOver = useCallback((e: DragEvent<HTMLElement>) => {
    e.preventDefault();
    e.stopPropagation();

    setState((prev) => ({
      ...prev,
      isDragOver: true,
    }));
  }, []);

  const onDragLeave = useCallback((e: DragEvent<HTMLElement>) => {
    e.preventDefault();
    e.stopPropagation();

    // Only reset if leaving the main container (not child elements)
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      setState((prev) => ({
        ...prev,
        isDragOver: false,
        isDragActive: false,
      }));
    }
  }, []);

  const onDrop = useCallback(
    (e: DragEvent<HTMLElement>) => {
      e.preventDefault();
      e.stopPropagation();

      setState({
        isDragOver: false,
        isDragActive: false,
      });

      const files = getFilesFromEvent(e);
      if (files.length > 0) {
        onFileDrop(files);
      }
    },
    [getFilesFromEvent, onFileDrop]
  );

  const reset = useCallback(() => {
    setState({
      isDragOver: false,
      isDragActive: false,
    });
  }, []);

  return {
    ...state,
    onDragEnter,
    onDragOver,
    onDragLeave,
    onDrop,
    reset,
  };
};
