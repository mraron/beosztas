using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.IO;
using System.Diagnostics;
using PdfSharp.Pdf;
using PdfSharp.Drawing;
using PdfSharp.Drawing.Layout;
using System.Net.Http;
using System.Net;
using System.Xml;

namespace ID
{
    static class Program
    {
        static readonly HttpClient client = new HttpClient();
        static void Main()
        {
            var request = (HttpWebRequest)WebRequest.Create("http://fmg.hu/osztalyok");
            var response = (HttpWebResponse)request.GetResponse();
            var responseString = new StreamReader(response.GetResponseStream()).ReadToEnd();
            string substr = "", link = "", title = "";
            List<string> links = new List<string>(), HTMLs = new List<string>(), classes = new List<string>(), titles = new List<string>();
            bool linkS = false, titleP = false, titleS = false;
            for(int i = 0; i < responseString.Length; i++)
            {
                if (substr.Length < 6) substr += responseString[i];
                else substr = substr.Substring(1) + responseString[i];
                if (linkS)
                {
                    if (responseString[i] == '\"')
                    {
                        links.Add(link);
                        linkS = false;
                        titleP = true;
                        link = "";
                    }
                    else link += responseString[i];
                }
                if (titleS)
                {
                    if (responseString[i] == '<')
                    {
                        titles.Add(title);
                        titleS = false;
                        title = "";
                    }
                    else title += responseString[i];
                }
                if(titleP && responseString[i] == '>')
                {
                    titleS = true;
                    titleP = false;
                }
                if (substr == "href=\"") linkS = true;
            }
            bool BefT = true;
            int pastOs = 0;
            for (int i = 0; i < links.Count; i++)
            {
                if (links[i] == "/osztalyok") pastOs++;
                if (links[i] == "/tanarok") BefT = false;
                if (pastOs >= 2 && BefT && links[i].Length >= 3 && links[i].Substring(0, 3) == "/20")
                {
                    HTMLs.Add("http://fmg.hu" + links[i]);
                    classes.Add(titles[i]);
                }
            }
            List<string[]> Names = new List<string[]>();
            string name = "";
            for (int i = 0; i < HTMLs.Count; i++)
            {
                request = (HttpWebRequest)WebRequest.Create(HTMLs[i]);
                response = (HttpWebResponse)request.GetResponse();
                responseString = new StreamReader(response.GetResponseStream()).ReadToEnd();
                bool PS = false;
                substr = "";
                for (int j = 0; j < responseString.Length; j++)
                {
                    if (substr.Length < 3) substr += responseString[j];
                    else substr = substr.Substring(1) + responseString[j];
                    if (PS)
                    {
                        if (substr + responseString[j + 1] == "</p>")
                        {
                            name = name.Substring(0, name.Length - 2);
                            if(!name.Contains("Osztálykép"))
                            {
                                Names.Add(name.Split(new string[] { "<br />\n" }, StringSplitOptions.None));
                            }
                            PS = false;
                            name = "";
                        }
                        else name += responseString[j];
                    }
                    if (substr == "<p>") PS = true;
                }
            }
            for (int i = 0; i < Names.Count; i++)
            {
                Console.WriteLine(classes[i]);
                for (int j = 0; j < Names[i].Length; j++)
                {
                    Console.WriteLine(Names[i][j]);
                }
            }
            string[] files = Directory.GetFiles("Names", "*.txt");
            var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
            var random = new Random();
            for(int k = 0; k < Names.Count; k++)
            {
                int n = Names[k].Length;
                string[] passwords = new string[n];
                string fnPdf = "Names\\" + classes[k];
                for (int i = 0; i < n; i++)
                {
                    for (int j = 0; j < 10; j++)
                    {
                        passwords[i] += chars[random.Next(chars.Length)];
                    }
                    Console.WriteLine(passwords[i]);
                }
                using (var fileStream = new FileStream("IDs\\ID_" + fnPdf.Split('\\')[fnPdf.Split('\\').Length - 1] + ".txt", FileMode.OpenOrCreate))
                using (var streamWriter = new StreamWriter(fileStream))
                {
                    for (int i = 0; i < n; i++)
                    {
                        streamWriter.WriteLine(fnPdf.Split('\\')[fnPdf.Split('\\').Length - 1] + "#" + Names[k][i] + "#" + passwords[i]);
                    }
                    streamWriter.Flush();
                    streamWriter.Close();
                }
                PdfDocument document = new PdfDocument();
                document.Info.Title = fnPdf;
                List<PdfPage> pages = new List<PdfPage>();
                List<XGraphics> gfx = new List<XGraphics>();
                List<XTextFormatter> tf = new List<XTextFormatter>();
                pages.Add(document.AddPage());
                gfx.Add(XGraphics.FromPdfPage(pages[0]));
                tf.Add(new XTextFormatter(gfx[0]));
                XFont titleFont = new XFont("Times New Roman", 32, XFontStyle.Bold);
                XFont font = new XFont("Times New Roman", 20, XFontStyle.Regular);
                tf[0].Alignment = XParagraphAlignment.Center;
                tf[0].DrawString(fnPdf.Split('\\')[fnPdf.Split('\\').Length - 1], titleFont, XBrushes.Black, new XRect(10, 10, pages[0].Width - 20, pages[0].Height - 40), XStringFormats.TopLeft);
                XPen pen = new XPen(XColors.Black, 1);
                gfx[0].DrawLine(pen, 0, 70, pages[0].Width, 70);
                int pageIndex = 0, newline = 80;
                for (int i = 0; i < n; i++)
                {
                    tf[pageIndex].Alignment = XParagraphAlignment.Left;
                    tf[pageIndex].DrawString(Names[k][i], font, XBrushes.Black, new XRect(20, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                    tf[pageIndex].Alignment = XParagraphAlignment.Left;
                    tf[pageIndex].DrawString(passwords[i], font, XBrushes.Black, new XRect(pages[pageIndex].Width / 2 + 120, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                    gfx[pageIndex].DrawLine(pen, 0, newline - 10, pages[pageIndex].Width, newline - 10);
                    if (newline + 40 > pages[pageIndex].Height)
                    {
                        gfx[pageIndex].DrawLine(pen, 0, newline + 30, pages[pageIndex].Width, newline + 30);
                        pageIndex++;
                        pages.Add(document.AddPage());
                        gfx.Add(XGraphics.FromPdfPage(pages[pageIndex]));
                        tf.Add(new XTextFormatter(gfx[pageIndex]));
                        tf[pageIndex].Alignment = XParagraphAlignment.Center;
                        tf[pageIndex].DrawString(fnPdf.Split('\\')[fnPdf.Split('\\').Length - 1], titleFont, XBrushes.Black, new XRect(10, 10, pages[0].Width - 20, pages[0].Height - 10), XStringFormats.TopLeft);
                        gfx[pageIndex].DrawLine(pen, 0, 70, pages[pageIndex].Width, 70);
                        gfx[pageIndex - 1].DrawLine(pen, pages[pageIndex - 1].Width / 2 + 100, 70, pages[pageIndex - 1].Width / 2 + 100, newline - 10);
                        newline = 80;
                        tf[pageIndex].Alignment = XParagraphAlignment.Left;
                        tf[pageIndex].DrawString(Names[k][i], font, XBrushes.Black, new XRect(20, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                        tf[pageIndex].Alignment = XParagraphAlignment.Left;
                        tf[pageIndex].DrawString(passwords[i], font, XBrushes.Black, new XRect(pages[pageIndex].Width / 2 + 120, newline, pages[pageIndex].Width - 20, pages[pageIndex].Height - 40), XStringFormats.TopLeft);
                        gfx[pageIndex].DrawLine(pen, 0, newline - 10, pages[pageIndex].Width, newline - 10);
                    }
                    newline += 40;
                }
                gfx[pageIndex].DrawLine(pen, 0, newline - 10, pages[pageIndex].Width, newline - 10);
                gfx[pageIndex].DrawLine(pen, pages[pageIndex].Width / 2 + 100, 70, pages[pageIndex].Width / 2 + 100, newline - 10);
                string filename = fnPdf + ".pdf";
                document.Save(filename);
                Process.Start(filename);
            }
        }
    }
}