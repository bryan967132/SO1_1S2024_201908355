const express = require('express');
const app = express();
const bodyParser = require('body-parser');
const mongoose = require('mongoose');

// Connect to MongoDB
mongoose.connect('mongodb://database/mydatabase', { useNewUrlParser: true, useUnifiedTopology: true });
const db = mongoose.connection;
db.on('error', console.error.bind(console, 'MongoDB connection error:'));

// Define a schema
const PhotoSchema = new mongoose.Schema({
    base64: String,
    uploadDate: { type: Date, default: Date.now }
});
const PhotoModel = mongoose.model('Photo', PhotoSchema);

app.use(bodyParser.json());

app.post('/photos', async (req, res) => {
    try {
        const { base64 } = req.body;
        const newPhoto = new PhotoModel({ base64 });
        await newPhoto.save();
        res.status(201).json(newPhoto);
    } catch (error) {
        console.error(error);
        res.status(500).send('Server Error');
    }
});

app.listen(3000, () => {
    console.log('API server is running on port 3000');
});