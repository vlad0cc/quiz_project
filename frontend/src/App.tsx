import { Navigate, Route, Routes } from 'react-router-dom';

import { QuestionPage } from './pages/QuestionPage';
import { ResultPage } from './pages/ResultPage';
import { WelcomePage } from './pages/WelcomePage';

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<WelcomePage />} />
      <Route path="/quiz" element={<QuestionPage />} />
      <Route path="/result/:sessionID" element={<ResultPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
