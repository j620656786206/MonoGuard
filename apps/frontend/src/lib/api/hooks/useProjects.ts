import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  ProjectsService, 
  type ProjectListParams, 
  type CreateProjectPayload,
  type UpdateProjectPayload,
  type AnalyzeProjectOptions
} from '../services/projects';
import { Project } from '@monoguard/shared-types';

/**
 * Query keys for projects
 */
export const projectQueryKeys = {
  all: ['projects'] as const,
  lists: () => [...projectQueryKeys.all, 'list'] as const,
  list: (params?: ProjectListParams) => [...projectQueryKeys.lists(), params] as const,
  details: () => [...projectQueryKeys.all, 'detail'] as const,
  detail: (id: string) => [...projectQueryKeys.details(), id] as const,
  dependencies: (id: string) => [...projectQueryKeys.detail(id), 'dependencies'] as const,
  architecture: (id: string) => [...projectQueryKeys.detail(id), 'architecture'] as const,
  health: (id: string) => [...projectQueryKeys.detail(id), 'health'] as const,
};

/**
 * Hook to fetch projects with optional filtering
 */
export function useProjects(params?: ProjectListParams) {
  return useQuery({
    queryKey: projectQueryKeys.list(params),
    queryFn: () => ProjectsService.getProjects(params),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

/**
 * Hook to fetch a single project
 */
export function useProject(id: string, enabled = true) {
  return useQuery({
    queryKey: projectQueryKeys.detail(id),
    queryFn: () => ProjectsService.getProject(id),
    enabled: enabled && !!id,
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
}

/**
 * Hook to fetch project dependencies
 */
export function useProjectDependencies(id: string, enabled = true) {
  return useQuery({
    queryKey: projectQueryKeys.dependencies(id),
    queryFn: () => ProjectsService.getProjectDependencies(id),
    enabled: enabled && !!id,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
}

/**
 * Hook to fetch project architecture
 */
export function useProjectArchitecture(id: string, enabled = true) {
  return useQuery({
    queryKey: projectQueryKeys.architecture(id),
    queryFn: () => ProjectsService.getProjectArchitecture(id),
    enabled: enabled && !!id,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
}

/**
 * Hook to fetch project health
 */
export function useProjectHealth(id: string, enabled = true) {
  return useQuery({
    queryKey: projectQueryKeys.health(id),
    queryFn: () => ProjectsService.getProjectHealth(id),
    enabled: enabled && !!id,
    staleTime: 1 * 60 * 1000, // 1 minute
  });
}

/**
 * Hook to create a new project
 */
export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateProjectPayload) => ProjectsService.createProject(payload),
    onSuccess: () => {
      // Invalidate projects list to refetch with new project
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.lists() });
    },
  });
}

/**
 * Hook to update a project
 */
export function useUpdateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateProjectPayload }) =>
      ProjectsService.updateProject(id, payload),
    onSuccess: (data, { id }) => {
      // Update the specific project in cache
      queryClient.setQueryData(projectQueryKeys.detail(id), data);
      // Invalidate projects list to refetch
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.lists() });
    },
  });
}

/**
 * Hook to delete a project
 */
export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => ProjectsService.deleteProject(id),
    onSuccess: (_, id) => {
      // Remove project from cache
      queryClient.removeQueries({ queryKey: projectQueryKeys.detail(id) });
      // Invalidate projects list
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.lists() });
    },
  });
}

/**
 * Hook to analyze a project
 */
export function useAnalyzeProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, options }: { id: string; options?: AnalyzeProjectOptions }) =>
      ProjectsService.analyzeProject(id, options),
    onSuccess: (_, { id }) => {
      // Invalidate project data to refetch with updated analysis
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.detail(id) });
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.health(id) });
    },
  });
}

/**
 * Hook to upload project files
 */
export function useUploadProjectFiles() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ 
      id, 
      files, 
      onUploadProgress 
    }: { 
      id: string; 
      files: File[]; 
      onUploadProgress?: (progressEvent: any) => void;
    }) => ProjectsService.uploadProjectFiles(id, files, onUploadProgress),
    onSuccess: (_, { id }) => {
      // Invalidate project data
      queryClient.invalidateQueries({ queryKey: projectQueryKeys.detail(id) });
    },
  });
}