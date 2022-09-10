import { useEffect, useState } from 'react';
import styled from "styled-components";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTrash, faThumbsUp } from '@fortawesome/free-solid-svg-icons'

const Container = styled.div`
    h1 {
        text-align: center;
    }
`;

const ImageContainer = styled.div`
    display: flex;
    border: 2px solid red;
    padding: 1em;
`;

const DeleteIconContainer = styled.div`
    width: 25px;
    border: none;
`;

const DeleteIcon = styled(FontAwesomeIcon)``;

const ThumbsUpIcon = styled(FontAwesomeIcon)``;

// const StyledDeleteIcon = styled(DeleteIcon)`
//     width: 15px;
// `;

const FileList = ({imagesInfo, isLoading}) => {
    console.log("PROPS imagesInfo", imagesInfo);

    
    function handleDelete(id) {
        console.log("CLICKING")
        const requestOptions = {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' }
        }

        fetch("http://localhost:8000/images/" + id, requestOptions)
            .then(response => {
                console.log("RESPONSE STATUS", response.status);
                return response.json();
            })
    }

    return (
        <Container>


            {isLoading && <div>Loading...</div>}

            {/* {Array.isArray(data) && data.length ? data.map(function (image) { */}
            {imagesInfo.length && imagesInfo ? imagesInfo.map(function (image) {
                return (
                    <ImageContainer>
                        {image.name}
                        {image.file_path}
                        <DeleteIconContainer onClick={() => handleDelete(image.id)}>
                            <DeleteIcon icon={faTrash} />
                        </DeleteIconContainer>
                        <ThumbsUpIcon icon={faThumbsUp} />
                    </ImageContainer>
                )
            }) : <div>Upload a file to see it here....</div>

            }

        </Container>
    )
}

// https://www.tutsmake.com/react-thumbnail-image-preview-before-upload-tutorial/
const ImagePreview = () => {
    return (
        <div>
            hi
        </div>
    );
}

export default FileList;