import type { AnswerResult, Question, Result } from '../types';

const BASE = '/api/v1';

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
    },
    ...init,
  });

  if (response.status === 204) {
    return null as T;
  }

  if (!response.ok) {
    let errorMessage = 'Request failed';
    try {
      const errorBody = (await response.json()) as { error?: string };
      if (errorBody.error) {
        errorMessage = errorBody.error;
      }
    } catch {
      errorMessage = response.statusText || errorMessage;
    }
    throw new Error(errorMessage);
  }

  return (await response.json()) as T;
}

export async function createSession(): Promise<{ session_id: string }> {
  return request<{ session_id: string }>('/sessions', {
    method: 'POST',
  });
}

export async function getNextQuestion(sessionId: string): Promise<Question | null> {
  return request<Question | null>(`/sessions/${sessionId}/question`);
}

export async function submitAnswer(
  sessionId: string,
  questionId: number,
  chosen: string,
): Promise<AnswerResult> {
  return request<AnswerResult>(`/sessions/${sessionId}/answers`, {
    method: 'POST',
    body: JSON.stringify({
      question_id: questionId,
      chosen,
    }),
  });
}

export async function getResult(sessionId: string): Promise<Result> {
  return request<Result>(`/sessions/${sessionId}/result`);
}
