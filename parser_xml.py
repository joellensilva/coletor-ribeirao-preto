import pandas as pd
# import numpy as np
import xml.etree.cElementTree as et

url = "http://www.abrasf.org.br/nfse.xsd"
# url = "http://www.ginfes.com.br/tipos"
arquivo = []
numero_nf = []
data_emissao = []
valor_servicos = []
base_calculo = []
valor_liquido_nfse = []
valor_iss_retido = []
discriminacao = []
municipio_prestacao_servico = []
cnpj_prestador = []
valor_deducoes = []

for i in range(2983, 2992):
    try:
        arq = f'nf-ribeirao-preto-{i}.xml'
        tree = et.parse(f'/home/joellen/estagio/coletor-ribeirao-preto/output/{arq}')
        root = tree.getroot()
        
        try:
            # Número da nf
            # print("1")
            item = root.find(f'.//{{{url}}}Numero')
            # print(item)
            # break
            if item != None:
                numero_nf.append(item.text)
            else:
                numero_nf.append("")

            # Data de emissão
            # print("2")
            item = root.find(f'.//{{{url}}}DataEmissao')
            if item != None:
                data_emissao.append(item.text)
            else:
                data_emissao.append("")

            # Valor de serviços
            # print("3")
            item = root.find(f'.//{{{url}}}ValorServicos')
            if item != None:
                valor_servicos.append(item.text)
            else:
                valor_servicos.append("")

            #  Base Cálculo
            # print("4")
            item = root.find(f'.//{{{url}}}BaseCalculo')
            if item != None:
                base_calculo.append(item.text)
            else:
                base_calculo.append("")

            # Valor líquido da nf
            # print("5")
            item = root.find(f'.//{{{url}}}ValorLiquidoNfse')
            if item != None:
                valor_liquido_nfse.append(item.text)
            else:
                valor_liquido_nfse.append("")
 
            # Valor ISS retido
            # print("6")
            item = root.find(f'.//{{{url}}}ValorIss')
            if item != None:
                valor_iss_retido.append(item.text)
            else:
                valor_iss_retido.append("")
 
            # Discriminação
            # print("7")
            item = root.find(f'.//{{{url}}}Discriminacao')
            if item != None:
                discriminacao.append(item.text)
            else:
                discriminacao.append("")
 
            # Município de prestação do serviço
            # print("8")
            item = root.find(f'.//{{{url}}}CodigoMunicipio')
            if item != None:
                municipio_prestacao_servico.append(item.text)
            else:
                municipio_prestacao_servico.append("")

            # CNPJ do prestador de serviço
            # print("9")
            item = root.find(f'.//{{{url}}}Cnpj')
            if item != None:
                cnpj_prestador.append(item.text)
            else:
                cnpj_prestador.append("")

            # Valor deduções
            # print("10")
            item = root.find(f'.//{{{url}}}ValorDeducoes')
            if item != None:
                valor_deducoes.append(item.text)
            else:
                valor_deducoes.append("")
            
            # Nome do arquivo
            arquivo.append(arq)
            print(f"OK: {arq}")
        except Exception as e:
            print(f"Erro parsing XML (step 2) -> {i}: {e}")
            data = pd.DataFrame(list(zip(arquivo, numero_nf, data_emissao, valor_servicos, base_calculo, valor_liquido_nfse, valor_iss_retido, discriminacao, municipio_prestacao_servico, cnpj_prestador, valor_deducoes)), 
                                columns=["arquivo", "numero_nf", "data_emissao", "valor_servicos", "base_calculo", "valor_liquido_nfse", "valor_iss_retido", "discriminacao", "municipio_prestacao_servico", "cnpj_prestador", "valor_deducoes"])
            # print(arquivo, numero_nf, data_emissao, valor_servicos, base_calculo, valor_liquido_nfse, valor_iss_retido, discriminacao, municipio_prestacao_servico, cnpj_prestador, valor_deducoes)
            data.to_csv("dados-extraidos-xml.csv", index=False)
            break
        
    except Exception as err:
        print(f"Erro parsing XML (step 1) -> {i}: {err}")
        # break
        pass
        
data = pd.DataFrame(list(zip(arquivo, numero_nf, data_emissao, valor_servicos, base_calculo, valor_liquido_nfse, valor_iss_retido, discriminacao, municipio_prestacao_servico, cnpj_prestador, valor_deducoes)), 
                            columns=["arquivo", "numero_nf", "data_emissao", "valor_servicos", "base_calculo", "valor_liquido_nfse", "valor_iss_retido", "discriminacao", "municipio_prestacao_servico", "cnpj_prestador", "valor_deducoes"])
data.to_csv("dados-extraidos-xml.csv", index=False) 