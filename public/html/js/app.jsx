const { useState, useMemo } = React;

const App = () => {
    const [validFor, setValidFor] = useState("60m")
    const [text, setText] = useState('');
    const [errMsg, setErrMsg] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [token, setToken] = useState('');
    const [isCoped, setIsCoped] = useState(false);

    const sendingData = useMemo(() => {
        return JSON.stringify({ url: text })
    }, [validFor, text])

    const onetimeLink = useMemo(() => {
        if (token === '') return ''
        return `${window.location.href.replace('#', '')}${token}`
    }, [token])

    const isReadyLink = useMemo(() => { return onetimeLink === '' ? false : true }, [onetimeLink])

    const handleClick = () => {

        if (isLoading) return
        if (text.length < 5) return
        setToken('')
        setErrMsg('')
        setIsCoped(false)
        setIsLoading(true)
        $.post({
            url: "/api/",
            data: sendingData,
            dataType: "json",
            success: function (response) {
                setToken(response)
                setIsLoading(false)
            },
            error: function (error) {
                error.responseJSON && setErrMsg(error.responseJSON.error)
                setToken('')
                setIsLoading(false)
            }
        })
    }

    //    console.log('href', window.location.href.replace('#', ''))
    return (
        <div>
            <div className="container " style={{ marginTop: "30px" }} >
                <div className="col-xs-8 col-xs-offset-1  pb-4 m-0" style={{ padding: "20px 160px" }}>
                    <div class="input-group" >
                        <img src="/icon.png" class="bg-transparent rounded " alt="icon" style={{ width: '50px', height: '50px' }}></img>

                        <p className="pl-4 h1 text-secondary">Кликер</p>
                    </div>
                    <p className="lead pt-2" style={{ marginLeft: "75px" }}>Благодаря короткой ссылке клиентам не придётся видеть длинные url-адреса.</p>
                    <div class="input-group" style={{ marginTop: '40px' }}>

                        <input
                            type="text"
                            value={text}
                            className="form-control form-control-lg"
                            placeholder="Сократить ссылку" aria-label="Сократить ссылку" aria-describedby="button-addon2"
                            onChange={(event) => { setText(event.target.value) }}
                            style={{ boxShadow: 'none' }}
                        />
                        <div class="input-group-append">
                            <button class="btn btn-primary ml-2" disabled={text.length < 5} id="button-addon2" onClick={handleClick} >Сократить</button>
                        </div>
                    </div>


                    {isLoading && <div class="d-flex justify-content-center mt-4">
                        <div class="spinner-border text-success" style={{ width: '3rem', height: '3rem' }} role="status">
                            <span class="sr-only">Loading...</span>
                        </div>
                    </div>}
                    {!isLoading && isReadyLink && <div class="input-group" style={{ marginTop: '70px' }}>
                        <input
                            type="text"
                            value={onetimeLink}
                            className={isCoped ? "form-control form-control-lg is-valid" : "form-control form-control-lg"}
                            style={{ boxShadow: 'none' }}
                            placeholder="Одноразовая ссылка"
                            aria-label="Одноразовая ссылка"
                            aria-describedby="button-addon2" />
                        <div class="input-group-append">
                            <CopyToClipboard onCopy={() => setIsCoped(true)} text={onetimeLink}>
                                <button class="btn btn-outline-secondary ml-2" type="button" id="button-addon2" >Копировать в Clipboard</button>
                            </CopyToClipboard>
                        </div>

                    </div>}
                    {errMsg &&
                        <div class="alert alert-danger mt-4" role="alert">
                            {errMsg}
                        </div>
                    }
                </div>
            </div >
            <div className="container" style={{ marginTop: "30px" }}>
                <div className="col-xs-10 col-xs-offset-2 pt-2 mb-2 ">


                </div>
            </div>
        </div >
    )
}

ReactDOM.render(<App />, document.getElementById('app'));

