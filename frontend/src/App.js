import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Books from './components/Books';

function App() {
  return (
    <Router>
      <div>
        <h1>Library Books</h1>
        <Routes>
          <Route path="/" element={<Books />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
