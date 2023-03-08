import '../styles/blog.scss'

import { useEffect, useState } from "react"

const updateParams = (key, value) => {
    if(window.history.pushState){
        let sParams = new URLSearchParams(window.location.search)
        sParams.set(key, value)
        let newURL = window.location.protocol + "//" + window.location.host + window.location.pathname + '?' + sParams.toString()
        window.history.pushState({path: newURL}, '', newURL)
    }
}

function Pagination(props) {
    const [offset, setOffset] = useState(0)
    const [pageSize,setPageSize] = useState(5)
    const maxPage = Math.ceil(props.numItems/props.pageSize)
    const [curPage, setCurPage] = useState(0)
    useEffect(()=>{
        updateParams('page', curPage);
        props.parentPage(curPage)
    }, [curPage, offset, pageSize])

    const next = () => {
        if(curPage < maxPage-1) {
            setCurPage(curPage+1);
            if((curPage+1)%5 == 0 && curPage > 1) {
                setOffset(offset + 1);
            }
        }
    }

    const prev = () => {
        if(curPage > 0) {
            setCurPage(curPage-1);
            if((curPage)%5 == 0 && curPage > 1) {
                setOffset(offset - 1);
            }
        }
    }

    return(
        <div className="pagination">
            <span className={curPage > 0 ? "navigator" : "navigator-hidden"} onClick={prev}>Prev</span>
            <span className="pages">{Array.from({length: props.pageSize}, (_,i) => {if (i+offset*props.pageSize < maxPage) return <a key={i+offset*props.pageSize} className={i+offset*props.pageSize == curPage ? "page-active" : "page"} page={i+offset*props.pageSize} onClick={()=>setCurPage(i+offset*props.pageSize)}>{i+offset*props.pageSize+1}</a>})}</span>
            <span className={curPage < maxPage-1 ? "navigator" : "navigator-hidden"} onClick={next}>Next</span>
        </div>
    )    
}

export default Pagination