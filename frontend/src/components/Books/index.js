import React, { useEffect, useState } from 'react';
import axios from 'axios';

const Books = () => {
  const [books, setBooks] = useState([]);
  const [editingBook, setEditingBook] = useState(null);
  const [updatedBook, setUpdatedBook] = useState({
    title: '',
    author: '',
    price: 0,
  });
  const [newBook, setNewBook] = useState({
    title: '',
    author: '',
    price: 0,
  });

  useEffect(() => {
    const fetchBooks = async () => {
      try {
        const response = await axios.get('http://localhost:8080/books');
        setBooks(response.data);
      } catch (error) {
        console.error("Error fetching books:", error);
      }
    };

    fetchBooks();
  }, []);

  const handleEditClick = (book) => {
    setEditingBook(book._id);
    setUpdatedBook({ title: book.title, author: book.author, price: book.price });
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setUpdatedBook((prev) => ({
      ...prev,
      [name]: name === 'price' ? parseFloat(value) : value,
    }));
  };

  const handleSaveClick = async () => {
    try {
      await axios.put(`http://localhost:8080/books/${editingBook}`, updatedBook);
      setBooks((prev) =>
        prev.map((book) =>
          book._id === editingBook ? { ...book, ...updatedBook } : book
        )
      );
      setEditingBook(null);
    } catch (error) {
      console.error("Error updating book:", error);
    }
  };

  const handleCancelClick = () => {
    setEditingBook(null);
  };

  const handleDeleteClick = async (bookId) => {
    try {
      await axios.delete(`http://localhost:8080/books/${bookId}`);
      setBooks((prev) => prev.filter((book) => book._id !== bookId));
    } catch (error) {
      console.error("Error deleting book:", error);
    }
  };

  const handleNewBookInputChange = (e) => {
    const { name, value } = e.target;
    setNewBook((prev) => ({
      ...prev,
      [name]: name === 'price' ? parseFloat(value) : value,
    }));
  };

  const handleAddBookClick = async () => {
    try {
      const response = await axios.post('http://localhost:8080/books', newBook);
      setBooks((prev) => [...prev, { ...newBook, _id: response.data.id }]);
      setNewBook({ title: '', author: '', price: 0 }); 
    } catch (error) {
      console.error("Error adding book:", error);
    }
  };

  return (
    <div>
      <h2>Books List</h2>

      {/* New Book Form */}
      <div>
        <h3>Add a New Book</h3>
        <input
          type="text"
          name="title"
          placeholder="Title"
          value={newBook.title}
          onChange={handleNewBookInputChange}
        />
        <input
          type="text"
          name="author"
          placeholder="Author"
          value={newBook.author}
          onChange={handleNewBookInputChange}
        />
        <input
          type="number"
          name="price"
          placeholder="Price"
          value={newBook.price}
          onChange={handleNewBookInputChange}
        />
        <button onClick={handleAddBookClick}>Add Book</button>
      </div>

      {/* Books List */}
      <ul>
        {books.map((book) => (
          <li key={book._id}>
            {editingBook === book._id ? (
              <div>
                <input
                  type="text"
                  name="title"
                  value={updatedBook.title}
                  onChange={handleInputChange}
                />
                <input
                  type="text"
                  name="author"
                  value={updatedBook.author}
                  onChange={handleInputChange}
                />
                <input
                  type="number"
                  name="price"
                  value={updatedBook.price}
                  onChange={handleInputChange}
                />
                <button onClick={handleSaveClick}>Save</button>
                <button onClick={handleCancelClick}>Cancel</button>
              </div>
            ) : (
              <div>
                <span className="book-info">{book.title} by {book.author} - ${book.price}</span>
                <button onClick={() => handleEditClick(book)}>Update</button>
                <button onClick={() => handleDeleteClick(book._id)}>Delete</button>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Books;
